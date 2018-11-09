// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gofrs/uuid/v3"
	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rsrv"
)

const cookieName = "GROOT_SRV"

type server struct {
	local bool
	srv   *rsrv.Server
	quit  chan int

	mu      sync.RWMutex
	cookies map[string]*http.Cookie
}

func newServer(local bool, dir string, mux *http.ServeMux) *server {
	app := &server{
		local:   local,
		srv:     rsrv.New(dir),
		quit:    make(chan int),
		cookies: make(map[string]*http.Cookie),
	}
	go app.run()

	mux.Handle("/", app.wrap(app.rootHandle))
	mux.HandleFunc("/ping", app.srv.Ping)
	mux.Handle("/root-file-upload", app.wrap(app.uploadHandle))
	mux.Handle("/root-file-open", app.wrap(app.openHandle))
	mux.Handle("/refresh", app.wrap(app.refreshHandle))
	mux.HandleFunc("/plot-h1", app.srv.PlotH1)
	mux.HandleFunc("/plot-h2", app.srv.PlotH2)
	mux.HandleFunc("/plot-s2", app.srv.PlotS2)
	mux.HandleFunc("/plot-branch", app.srv.PlotTree)

	return app
}

func (srv *server) run() {
	defer srv.srv.Shutdown()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	srv.gc()
	for {
		select {
		case <-ticker.C:
			srv.gc()
		case <-srv.quit:
			return
		}
	}
}

func (srv *server) Shutdown() {
	close(srv.quit)
}

func (srv *server) gc() {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	for name, cookie := range srv.cookies {
		now := time.Now()
		if now.After(cookie.Expires) {
			delete(srv.cookies, name)
			cookie.MaxAge = -1
		}
	}
}

func (srv *server) expired(cookie *http.Cookie) bool {
	now := time.Now()
	return now.After(cookie.Expires)
}

func (srv *server) setCookie(w http.ResponseWriter, r *http.Request) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	cookie, err := r.Cookie(cookieName)
	if err != nil && err != http.ErrNoCookie {
		return err
	}

	if cookie != nil {
		return nil
	}

	cookie = &http.Cookie{
		Name:    cookieName,
		Value:   uuid.Must(uuid.NewV4()).String(),
		Expires: time.Now().Add(24 * time.Hour),
	}
	srv.cookies[cookie.Value] = cookie
	http.SetCookie(w, cookie)
	return nil
}

func (srv *server) wrap(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := srv.setCookie(w, r)
		if err != nil {
			log.Printf("error retrieving cookie: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := fn(w, r); err != nil {
			log.Printf("error %q: %v\n", r.URL.Path, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (srv *server) rootHandle(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		// ok
	default:
		return fmt.Errorf("invalid request %q for /", r.Method)
	}

	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))

	t, err := template.New("upload").Parse(page)
	if err != nil {
		return err
	}

	srv.ping(r)

	return t.Execute(w, struct {
		Token string
		Local bool
	}{token, srv.local})
}

func (srv *server) uploadHandle(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return errors.Wrap(err, "could not retrieve cookie")
	}

	defer r.Body.Close()
	req, err := http.NewRequest(http.MethodPost, "/file-upload", r.Body)
	if err != nil {
		return errors.Wrap(err, "could not create upload-file request")
	}
	req.AddCookie(cookie)
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))

	ww := newResponseWriter()
	srv.srv.UploadFile(ww, req)

	if ww.code != http.StatusOK {
		w.WriteHeader(ww.code)
		return errors.Errorf("could not upload file")
	}

	nodes, err := srv.nodes(r)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(nodes)
}

func (srv *server) openHandle(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return errors.Wrap(err, "could not retrieve cookie")
	}

	err = r.ParseMultipartForm(500 << 20)
	if err != nil {
		return errors.Wrapf(err, "could not parse multipart form")
	}
	fname := r.PostFormValue("uri")
	if fname == "" {
		w.WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(w).Encode(nil)
	}

	body := new(bytes.Buffer)
	err = json.NewEncoder(body).Encode(rsrv.OpenFileRequest{URI: fname})
	if err != nil {
		return errors.Wrap(err, "could not encode open-file request")
	}

	req, err := http.NewRequest(http.MethodPost, "/file-open", body)
	if err != nil {
		return errors.Wrap(err, "could not create open-file request")
	}
	req.AddCookie(cookie)

	ww := newResponseWriter()
	srv.srv.OpenFile(ww, req)
	body.Truncate(0)

	if ww.code != http.StatusOK {
		w.WriteHeader(ww.code)
		return errors.Errorf("could not open file %q", fname)
	}

	nodes, err := srv.nodes(r)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(nodes)
}

func (srv *server) refreshHandle(w http.ResponseWriter, r *http.Request) error {
	nodes, err := srv.nodes(r)
	if err != nil {
		if err == http.ErrNoCookie {
			return json.NewEncoder(w).Encode(nil)
		}
		return err
	}

	return json.NewEncoder(w).Encode(nodes)
}

func (srv *server) nodes(r *http.Request) ([]jsNode, error) {
	db, err := srv.srv.DB(r)
	if err != nil {
		return nil, err
	}

	var nodes []jsNode
	uris := db.Files()
	for _, uri := range uris {
		err = db.Tx(uri, func(f *riofs.File) error {
			node, err := fileJsTree(f, uri)
			if err != nil {
				return err
			}
			nodes = append(nodes, node...)
			return nil
		})
		if err != nil {
			return nil, errors.Wrapf(err, "could not build nodes-tree for %q", uri)
		}
	}

	sort.Sort(jsNodes(nodes))
	return nodes, nil
}

func (srv *server) ping(r *http.Request) error {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		return err
	}
	req.AddCookie(cookie)

	ww := newResponseWriter()
	srv.srv.Ping(ww, req)

	if ww.code != http.StatusOK {
		return errors.Errorf("could not ping")
	}

	return nil
}
