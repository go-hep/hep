// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsrv

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	uuid "github.com/hashicorp/go-uuid"
)

const (
	cookieName = "GROOT_SRV"
)

// Server serves and manages ROOT files.
type Server struct {
	quit chan int

	mu       sync.RWMutex
	cookies  map[string]*http.Cookie
	sessions map[string]*DB

	dir string
}

// New creates a new server.
func New(dir string) *Server {
	srv := &Server{
		quit:     make(chan int),
		cookies:  make(map[string]*http.Cookie),
		sessions: make(map[string]*DB),
		dir:      dir,
	}

	go srv.run()
	return srv
}

// Shutdown shuts the server down.
func (srv *Server) Shutdown() {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	for name := range srv.cookies {
		srv.sessions[name].Close()
		delete(srv.sessions, name)
		delete(srv.cookies, name)
	}
	close(srv.quit)
}

func (srv *Server) run() {
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

func (srv *Server) gc() {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	for name, cookie := range srv.cookies {
		now := time.Now()
		if now.After(cookie.Expires) {
			srv.sessions[name].Close()
			delete(srv.sessions, name)
			delete(srv.cookies, name)
			cookie.MaxAge = -1
		}
	}
}

func (srv *Server) wrap(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
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

func (srv *Server) setCookie(w http.ResponseWriter, r *http.Request) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	cookie, err := r.Cookie(cookieName)
	if err != nil && err != http.ErrNoCookie {
		return err
	}

	if cookie != nil {
		if v, ok := srv.sessions[cookie.Value]; v == nil || !ok {
			srv.sessions[cookie.Value] = NewDB(filepath.Join(srv.dir, cookie.Value))
			srv.cookies[cookie.Value] = cookie
		}
		return nil
	}

	v, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("could not generate UUID: %w", err)
	}

	cookie = &http.Cookie{
		Name:    cookieName,
		Value:   v,
		Expires: time.Now().Add(24 * time.Hour),
	}
	srv.sessions[cookie.Value] = NewDB(filepath.Join(srv.dir, cookie.Value))
	srv.cookies[cookie.Value] = cookie
	http.SetCookie(w, cookie)
	return nil
}

func (srv *Server) cookie(r *http.Request) (*http.Cookie, error) {
	srv.mu.RLock()
	defer srv.mu.RUnlock()
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, err
	}

	if cookie == nil {
		return nil, http.ErrNoCookie
	}
	return srv.cookies[cookie.Value], nil
}

func (srv *Server) db(r *http.Request) (*DB, error) {
	srv.mu.RLock()
	defer srv.mu.RUnlock()
	cookie, err := srv.cookie(r)
	if err != nil {
		return nil, err
	}
	if cookie == nil {
		return nil, http.ErrNoCookie
	}
	return srv.sessions[cookie.Value], nil
}

// DB returns the underlying data base of files associated with the user
// identified by their cookie.
func (srv *Server) DB(r *http.Request) (*DB, error) {
	return srv.db(r)
}
