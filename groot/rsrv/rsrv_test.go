// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsrv

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
	"gonum.org/v1/plot/cmpimg"
)

var (
	srv *Server
)

func TestMain(m *testing.M) {
	dir, err := ioutil.TempDir("", "groot-rsrv-")
	if err != nil {
		log.Panicf("could not create temporary directory: %v", err)
	}
	defer os.RemoveAll(dir)

	srv = New(dir)
	setupCookie(srv)

	os.Exit(m.Run())
}

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/open-file", srv.OpenFile)
	mux.HandleFunc("/upload-file", srv.UploadFile)
	mux.HandleFunc("/close-file", srv.CloseFile)
	mux.HandleFunc("/list-files", srv.ListFiles)
	mux.HandleFunc("/list-dirs", srv.Dirent)
	mux.HandleFunc("/list-tree", srv.Tree)
	mux.HandleFunc("/plot-h1", srv.PlotH1)
	mux.HandleFunc("/plot-h2", srv.PlotH2)
	mux.HandleFunc("/plot-s2", srv.PlotS2)
	mux.HandleFunc("/plot-tree", srv.PlotTree)

	return httptest.NewServer(mux)
}

func TestOpenFile(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	local, err := filepath.Abs("../testdata/simple.root")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	for _, tc := range []struct {
		uri    string
		status int
	}{
		{"https://github.com/go-hep/hep/raw/master/groot/testdata/simple.root", http.StatusOK},
		{"root://ccxrootdgotest.in2p3.fr:9001/tmp/rootio/testdata/simple.root", http.StatusOK},
		{"file://" + local, http.StatusOK},
	} {
		t.Run(tc.uri, func(t *testing.T) {
			testOpenFile(t, ts, tc.uri, tc.status)
			defer testCloseFile(t, ts, tc.uri)
		})
	}
}

func TestDoubleOpenFile(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	testOpenFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/simple.root", 0)
	testOpenFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/simple.root", http.StatusConflict)
}

func testOpenFile(t *testing.T, ts *httptest.Server, uri string, status int) {
	t.Helper()

	req := OpenFileRequest{URI: uri}

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(req)
	if err != nil {
		t.Fatalf("could not encode request: %v", err)
	}

	hreq, err := http.NewRequest(http.MethodPost, ts.URL+"/open-file", body)
	if err != nil {
		t.Fatalf("could not create http request: %v", err)
	}
	srv.addCookies(hreq)

	hresp, err := ts.Client().Do(hreq)
	if err != nil {
		t.Fatalf("could not post http request: %v", err)
	}
	defer hresp.Body.Close()

	if got, want := hresp.StatusCode, status; got != want && want != 0 {
		t.Fatalf("invalid status code: got=%v, want=%v", got, want)
	}
}

func TestUploadFile(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	local, err := filepath.Abs("../testdata/simple.root")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	for _, tc := range []struct {
		dst, src string
		status   int
	}{
		{"foo.root", local, http.StatusOK},
	} {
		t.Run(tc.dst, func(t *testing.T) {
			testUploadFile(t, ts, tc.dst, tc.src, tc.status)
			defer testCloseFile(t, ts, tc.dst)
		})
	}
}

func testUploadFile(t *testing.T, ts *httptest.Server, dst, src string, status int) {
	t.Helper()

	body := new(bytes.Buffer)
	mpart := multipart.NewWriter(body)
	req, err := mpart.CreateFormField("groot-dst")
	if err != nil {
		t.Fatalf("could not create json-request form field: %v", err)
	}
	_, err = req.Write([]byte(dst))
	if err != nil {
		t.Fatalf("could not fill destination field: %v", err)
	}

	w, err := mpart.CreateFormFile("groot-file", src)
	if err != nil {
		t.Fatalf("could not create form-file: %v", err)
	}
	{
		f, err := os.Open(src)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer f.Close()

		_, err = io.CopyBuffer(w, f, make([]byte, 16*1024*1024))
		if err != nil {
			t.Fatalf("could not copy file: %v", err)
		}
	}

	if err := mpart.Close(); err != nil {
		t.Fatalf("could not close multipart form data: %v", err)
	}

	hreq, err := http.NewRequest(http.MethodPost, ts.URL+"/upload-file", body)
	if err != nil {
		t.Fatalf("could not create http request: %v", err)
	}
	srv.addCookies(hreq)
	hreq.Header.Set("Content-Type", mpart.FormDataContentType())

	hresp, err := ts.Client().Do(hreq)
	if err != nil {
		t.Fatalf("could not post http request: %v", err)
	}
	defer hresp.Body.Close()

	if got, want := hresp.StatusCode, status; got != want && want != 0 {
		t.Fatalf("invalid status code: got=%v, want=%v", got, want)
	}
}

func TestCloseFile(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	testOpenFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/simple.root", 0)
	testCloseFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/simple.root")
	testOpenFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/simple.root", http.StatusOK)
	testCloseFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/simple.root")
}

func testCloseFile(t *testing.T, ts *httptest.Server, uri string) {
	t.Helper()

	req := CloseFileRequest{URI: uri}
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(req)
	if err != nil {
		t.Fatalf("could not encode request: %v", err)
	}

	hreq, err := http.NewRequest(http.MethodPost, ts.URL+"/close-file", body)
	if err != nil {
		t.Fatalf("could not create http request: %v", err)
	}
	srv.addCookies(hreq)

	hresp, err := ts.Client().Do(hreq)
	if err != nil {
		t.Fatalf("could not post http request: %v", err)
	}
	defer hresp.Body.Close()

	if hresp.StatusCode != http.StatusOK {
		t.Fatalf("could not close file %q: %v", uri, hresp.StatusCode)
	}
}

func TestListFiles(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	testOpenFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/simple.root", 0)
	testOpenFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/dirs-6.14.00.root", http.StatusOK)
	testListFiles(t, ts, []File{
		{"https://github.com/go-hep/hep/raw/master/groot/testdata/simple.root", 60600},
		{"https://github.com/go-hep/hep/raw/master/groot/testdata/dirs-6.14.00.root", 61400},
	})
	testCloseFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/simple.root")
	testCloseFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/dirs-6.14.00.root")
}

func testListFiles(t *testing.T, ts *httptest.Server, want []File) {
	t.Helper()

	hreq, err := http.NewRequest(http.MethodPost, ts.URL+"/list-files", nil)
	if err != nil {
		t.Fatalf("could not create http request: %v", err)
	}
	srv.addCookies(hreq)

	hresp, err := ts.Client().Do(hreq)
	if err != nil {
		t.Fatalf("could not post http request: %v", err)
	}
	defer hresp.Body.Close()

	if hresp.StatusCode != http.StatusOK {
		t.Fatalf("could not list files: %v", hresp.StatusCode)
	}

	var resp ListResponse
	err = json.NewDecoder(hresp.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	got := resp.Files
	sort.Slice(got, func(i, j int) bool {
		return got[i].URI < got[j].URI
	})
	sort.Slice(want, func(i, j int) bool {
		return want[i].URI < want[j].URI
	})

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid ls content:\ngot= %v\nwant=%v\n", got, want)
	}
}

func TestDirent(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	testOpenFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/simple.root", http.StatusOK)
	defer testCloseFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/simple.root")

	testDirent(t, ts, DirentRequest{
		URI:       "https://github.com/go-hep/hep/raw/master/groot/testdata/simple.root",
		Dir:       "/",
		Recursive: false,
	}, []string{
		"/",
		"/tree",
	})

	testOpenFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/dirs-6.14.00.root", http.StatusOK)
	defer testCloseFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/dirs-6.14.00.root")

	testDirent(t, ts, DirentRequest{
		URI:       "https://github.com/go-hep/hep/raw/master/groot/testdata/dirs-6.14.00.root",
		Dir:       "/",
		Recursive: false,
	}, []string{
		"/",
		"/dir1",
		"/dir2",
		"/dir3",
	})
	testDirent(t, ts, DirentRequest{
		URI:       "https://github.com/go-hep/hep/raw/master/groot/testdata/dirs-6.14.00.root",
		Dir:       "/",
		Recursive: true,
	}, []string{
		"/",
		"/dir1",
		"/dir1/dir11",
		"/dir1/dir11/h1",
		"/dir2",
		"/dir3",
	})
	testDirent(t, ts, DirentRequest{
		URI:       "https://github.com/go-hep/hep/raw/master/groot/testdata/dirs-6.14.00.root",
		Dir:       "/dir1",
		Recursive: false,
	}, []string{
		"/dir1",
		"/dir1/dir11",
	})
	testDirent(t, ts, DirentRequest{
		URI:       "https://github.com/go-hep/hep/raw/master/groot/testdata/dirs-6.14.00.root",
		Dir:       "/dir1",
		Recursive: true,
	}, []string{
		"/dir1",
		"/dir1/dir11",
		"/dir1/dir11/h1",
	})
}

func testDirent(t *testing.T, ts *httptest.Server, req DirentRequest, content []string) {
	t.Helper()

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(req)
	if err != nil {
		t.Fatalf("could not encode request: %v", err)
	}

	hreq, err := http.NewRequest(http.MethodPost, ts.URL+"/list-dirs", body)
	if err != nil {
		t.Fatalf("could not create http request: %v", err)
	}
	srv.addCookies(hreq)

	hresp, err := ts.Client().Do(hreq)
	if err != nil {
		t.Fatalf("could not post http request: %v", err)
	}
	defer hresp.Body.Close()

	if hresp.StatusCode != http.StatusOK {
		t.Fatalf("could not list dirs: %v", hresp.StatusCode)
	}

	var resp DirentResponse
	err = json.NewDecoder(hresp.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	var got []string
	for _, f := range resp.Content {
		got = append(got, f.Path)
	}

	sort.Strings(got)
	sort.Strings(content)

	if !reflect.DeepEqual(got, content) {
		t.Fatalf("invalid dirent content: (req=%#v)\ngot= %v\nwant=%v\n", req, got, content)
	}
}

func TestTree(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	const uri = "https://github.com/go-hep/hep/raw/master/groot/testdata/small-flat-tree.root"
	testOpenFile(t, ts, uri, http.StatusOK)
	defer testCloseFile(t, ts, uri)

	for _, tc := range []struct {
		req  TreeRequest
		want Tree
	}{
		{
			req: TreeRequest{
				URI: uri,
				Obj: "tree",
			},
			want: Tree{
				Type:    "TTree",
				Name:    "tree",
				Title:   "my tree title",
				Entries: 100,
				Branches: []Branch{
					{Type: "TBranch", Name: "Int32", Leaves: []Leaf{{Type: "int32", Name: "Int32"}}},
					{Type: "TBranch", Name: "Int64", Leaves: []Leaf{{Type: "int64", Name: "Int64"}}},
					{Type: "TBranch", Name: "UInt32", Leaves: []Leaf{{Type: "uint32", Name: "UInt32"}}},
					{Type: "TBranch", Name: "UInt64", Leaves: []Leaf{{Type: "uint64", Name: "UInt64"}}},
					{Type: "TBranch", Name: "Float32", Leaves: []Leaf{{Type: "float32", Name: "Float32"}}},
					{Type: "TBranch", Name: "Float64", Leaves: []Leaf{{Type: "float64", Name: "Float64"}}},
					{Type: "TBranch", Name: "Str", Leaves: []Leaf{{Type: "string", Name: "Str"}}},
					{Type: "TBranch", Name: "ArrayInt32", Leaves: []Leaf{{Type: "int32", Name: "ArrayInt32"}}},
					{Type: "TBranch", Name: "ArrayInt64", Leaves: []Leaf{{Type: "int64", Name: "ArrayInt64"}}},
					{Type: "TBranch", Name: "ArrayUInt32", Leaves: []Leaf{{Type: "uint32", Name: "ArrayInt32"}}}, // FIXME(sbinet): ref-file had a typo (should be ArrayUInt32)
					{Type: "TBranch", Name: "ArrayUInt64", Leaves: []Leaf{{Type: "uint64", Name: "ArrayInt64"}}}, // FIXME(sbinet): ref-file had a typo (should be ArrayUInt64)
					{Type: "TBranch", Name: "ArrayFloat32", Leaves: []Leaf{{Type: "float32", Name: "ArrayFloat32"}}},
					{Type: "TBranch", Name: "ArrayFloat64", Leaves: []Leaf{{Type: "float64", Name: "ArrayFloat64"}}},
					{Type: "TBranch", Name: "N", Leaves: []Leaf{{Type: "int32", Name: "N"}}},
					{Type: "TBranch", Name: "SliceInt32", Leaves: []Leaf{{Type: "int32", Name: "SliceInt32"}}},
					{Type: "TBranch", Name: "SliceInt64", Leaves: []Leaf{{Type: "int64", Name: "SliceInt64"}}},
					{Type: "TBranch", Name: "SliceUInt32", Leaves: []Leaf{{Type: "uint32", Name: "SliceInt32"}}}, // FIXME(sbinet): ref-file had a typo (should be SliceUInt32)
					{Type: "TBranch", Name: "SliceUInt64", Leaves: []Leaf{{Type: "uint64", Name: "SliceInt64"}}}, // FIXME(sbinet): ref-file had a typo (should be SliceUInt64)
					{Type: "TBranch", Name: "SliceFloat32", Leaves: []Leaf{{Type: "float32", Name: "SliceFloat32"}}},
					{Type: "TBranch", Name: "SliceFloat64", Leaves: []Leaf{{Type: "float64", Name: "SliceFloat64"}}},
				},
				Leaves: []Leaf{
					{Type: "int32", Name: "Int32"},
					{Type: "int64", Name: "Int64"},
					{Type: "uint32", Name: "UInt32"},
					{Type: "uint64", Name: "UInt64"},
					{Type: "float32", Name: "Float32"},
					{Type: "float64", Name: "Float64"},
					{Type: "string", Name: "Str"},
					{Type: "int32", Name: "ArrayInt32"},
					{Type: "int64", Name: "ArrayInt64"},
					{Type: "uint32", Name: "ArrayInt32"}, // FIXME(sbinet): ref-file had a typo (should be ArrayUInt32)
					{Type: "uint64", Name: "ArrayInt64"}, // FIXME(sbinet): ref-file had a typo (should be ArrayUInt64)
					{Type: "float32", Name: "ArrayFloat32"},
					{Type: "float64", Name: "ArrayFloat64"},
					{Type: "int32", Name: "N"},
					{Type: "int32", Name: "SliceInt32"},
					{Type: "int64", Name: "SliceInt64"},
					{Type: "uint32", Name: "SliceInt32"}, // FIXME(sbinet): ref-file had a typo (should be SliceUInt32)
					{Type: "uint64", Name: "SliceInt64"}, // FIXME(sbinet): ref-file had a typo (should be SliceUInt64)
					{Type: "float32", Name: "SliceFloat32"},
					{Type: "float64", Name: "SliceFloat64"},
				},
			},
		},
	} {
		t.Run(tc.want.Name, func(t *testing.T) {
			var resp TreeResponse
			testTree(t, ts, tc.req, &resp)

			if !reflect.DeepEqual(resp.Tree, tc.want) {
				t.Fatalf("invalid tree:\ngot= %#v\nwant=%#v", resp.Tree, tc.want)
			}
		})
	}
}

func testTree(t *testing.T, ts *httptest.Server, req TreeRequest, resp *TreeResponse) {
	t.Helper()

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(req)
	if err != nil {
		t.Fatalf("could not encode request: %v", err)
	}

	hreq, err := http.NewRequest(http.MethodPost, ts.URL+"/list-tree", body)
	if err != nil {
		t.Fatalf("could not create http request: %v", err)
	}
	srv.addCookies(hreq)

	hresp, err := ts.Client().Do(hreq)
	if err != nil {
		t.Fatalf("could not post http request: %v", err)
	}
	defer hresp.Body.Close()

	if hresp.StatusCode != http.StatusOK {
		t.Fatalf("could not plot h1: %v", hresp.StatusCode)
	}

	err = json.NewDecoder(hresp.Body).Decode(resp)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
}

func TestPlotH1(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	testOpenFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/dirs-6.14.00.root", http.StatusOK)
	defer testCloseFile(t, ts, "https://github.com/go-hep/hep/raw/master/groot/testdata/dirs-6.14.00.root")

	const uri = "https://github.com/go-hep/hep/raw/master/hbook/rootcnv/testdata/gauss-h1.root"
	testOpenFile(t, ts, uri, http.StatusOK)
	defer testCloseFile(t, ts, uri)

	for _, tc := range []struct {
		req  PlotH1Request
		want string
	}{
		{
			req: PlotH1Request{
				URI: uri,
				Obj: "h1f",
			},
			want: "testdata/h1f_golden.png",
		},
		{
			req: PlotH1Request{
				URI: "https://github.com/go-hep/hep/raw/master/groot/testdata/dirs-6.14.00.root",
				Dir: "/dir1/dir11",
				Obj: "h1",
			},
			want: "testdata/h1_golden.png",
		},
		{
			req: PlotH1Request{
				URI: uri,
				Obj: "h1f",
				Options: PlotOptions{
					Type:      "png",
					Title:     "My Title",
					X:         "X axis [GeV]",
					Y:         "Y axis [A.U]",
					FillColor: color.RGBA{0, 0, 255, 255},
					Line: LineStyle{
						Color: color.Black,
					},
				},
			},
			want: "testdata/h1f_options_golden.png",
		},
		{
			req: PlotH1Request{
				URI: "https://github.com/go-hep/hep/raw/master/groot/testdata/dirs-6.14.00.root",
				Dir: "/dir1/dir11",
				Obj: "h1",
				Options: PlotOptions{
					Type: "pdf",
				},
			},
			want: "testdata/h1_golden.pdf",
		},
		{
			req: PlotH1Request{
				URI: "https://github.com/go-hep/hep/raw/master/groot/testdata/dirs-6.14.00.root",
				Dir: "/dir1/dir11",
				Obj: "h1",
				Options: PlotOptions{
					Type: "svg",
				},
			},
			want: "testdata/h1_golden.svg",
		},
	} {
		t.Run(tc.want, func(t *testing.T) {
			var resp PlotResponse
			testPlotH1(t, ts, tc.req, &resp)

			raw, err := base64.StdEncoding.DecodeString(resp.Data)
			if err != nil {
				t.Fatal(err)
			}

			if *cmpimg.GenerateTestData {
				ioutil.WriteFile(tc.want, raw, 0644)
			}

			want, err := ioutil.ReadFile(tc.want)
			if err != nil {
				t.Fatal(err)
			}

			typ := tc.req.Options.Type
			if typ == "" {
				typ = "png"
			}
			if ok, err := cmpimg.Equal(typ, raw, want); !ok || err != nil {
				ioutil.WriteFile(strings.Replace(tc.want, "_golden", "", -1), raw, 0644)
				t.Fatalf("reference files differ: err=%v ok=%v", err, ok)
			}
		})
	}
}

func testPlotH1(t *testing.T, ts *httptest.Server, req PlotH1Request, resp *PlotResponse) {
	t.Helper()

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(req)
	if err != nil {
		t.Fatalf("could not encode request: %v", err)
	}

	hreq, err := http.NewRequest(http.MethodPost, ts.URL+"/plot-h1", body)
	if err != nil {
		t.Fatalf("could not create http request: %v", err)
	}
	srv.addCookies(hreq)

	hresp, err := ts.Client().Do(hreq)
	if err != nil {
		t.Fatalf("could not post http request: %v", err)
	}
	defer hresp.Body.Close()

	if hresp.StatusCode != http.StatusOK {
		t.Fatalf("could not plot h1: %v", hresp.StatusCode)
	}

	err = json.NewDecoder(hresp.Body).Decode(resp)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
}

func TestPlotH2(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	const uri = "https://github.com/go-hep/hep/raw/master/hbook/rootcnv/testdata/gauss-h2.root"
	testOpenFile(t, ts, uri, http.StatusOK)
	defer testCloseFile(t, ts, uri)

	for _, tc := range []struct {
		req  PlotH2Request
		want string
	}{
		{
			req: PlotH2Request{
				URI: uri,
				Obj: "h2f",
			},
			want: "testdata/h2f_golden.png",
		},
		{
			req: PlotH2Request{
				URI: uri,
				Dir: "/",
				Obj: "h2d",
				Options: PlotOptions{
					Type: "png",
				},
			},
			want: "testdata/h2d_golden.png",
		},
		{
			req: PlotH2Request{
				URI: uri,
				Dir: "/",
				Obj: "h2d",
				Options: PlotOptions{
					Type:  "png",
					Title: "My Title",
					X:     "X-axis [GeV]",
					Y:     "Y-axis [GeV]",
				},
			},
			want: "testdata/h2d_options_golden.png",
		},
		{
			req: PlotH2Request{
				URI: uri,
				Dir: "/",
				Obj: "h2d",
				Options: PlotOptions{
					Type: "pdf",
				},
			},
			want: "testdata/h2d_golden.pdf",
		},
		{
			req: PlotH2Request{
				URI: uri,
				Dir: "/",
				Obj: "h2d",
				Options: PlotOptions{
					Type: "svg",
				},
			},
			want: "testdata/h2d_golden.svg",
		},
	} {
		t.Run(tc.want, func(t *testing.T) {
			var resp PlotResponse
			testPlotH2(t, ts, tc.req, &resp)

			raw, err := base64.StdEncoding.DecodeString(resp.Data)
			if err != nil {
				t.Fatal(err)
			}

			if *cmpimg.GenerateTestData {
				ioutil.WriteFile(tc.want, raw, 0644)
			}

			want, err := ioutil.ReadFile(tc.want)
			if err != nil {
				t.Fatal(err)
			}

			typ := tc.req.Options.Type
			if typ == "" {
				typ = "png"
			}
			if ok, err := cmpimg.Equal(typ, raw, want); !ok || err != nil {
				ioutil.WriteFile(strings.Replace(tc.want, "_golden", "", -1), raw, 0644)
				t.Fatalf("reference files differ: err=%v ok=%v", err, ok)
			}
		})
	}
}

func testPlotH2(t *testing.T, ts *httptest.Server, req PlotH2Request, resp *PlotResponse) {
	t.Helper()

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(req)
	if err != nil {
		t.Fatalf("could not encode request: %v", err)
	}

	hreq, err := http.NewRequest(http.MethodPost, ts.URL+"/plot-h2", body)
	if err != nil {
		t.Fatalf("could not create http request: %v", err)
	}
	srv.addCookies(hreq)

	hresp, err := ts.Client().Do(hreq)
	if err != nil {
		t.Fatalf("could not post http request: %v", err)
	}
	defer hresp.Body.Close()

	if hresp.StatusCode != http.StatusOK {
		t.Fatalf("could not plot h1: %v", hresp.StatusCode)
	}

	err = json.NewDecoder(hresp.Body).Decode(resp)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
}

func TestPlotS2(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	const uri = "https://github.com/go-hep/hep/raw/master/groot/testdata/graphs.root"
	testOpenFile(t, ts, uri, http.StatusOK)
	defer testCloseFile(t, ts, uri)

	for _, tc := range []struct {
		req  PlotS2Request
		want string
	}{
		{
			req: PlotS2Request{
				URI: uri,
				Obj: "tg",
			},
			want: "testdata/tg_golden.png",
		},
		{
			req: PlotS2Request{
				URI: uri,
				Dir: "/",
				Obj: "tge",
				Options: PlotOptions{
					Type: "png",
				},
			},
			want: "testdata/tge_golden.png",
		},
		{
			req: PlotS2Request{
				URI: uri,
				Dir: "/",
				Obj: "tgae",
				Options: PlotOptions{
					Type:  "png",
					Title: "My Title",
					X:     "X-axis [GeV]",
					Y:     "Y-axis [GeV]",
					Line: LineStyle{
						Color: color.RGBA{B: 255, A: 255},
					},
				},
			},
			want: "testdata/tgae_options_golden.png",
		},
		{
			req: PlotS2Request{
				URI: uri,
				Dir: "/",
				Obj: "tgae",
				Options: PlotOptions{
					Type: "pdf",
				},
			},
			want: "testdata/tgae_golden.pdf",
		},
		{
			req: PlotS2Request{
				URI: uri,
				Dir: "/",
				Obj: "tgae",
				Options: PlotOptions{
					Type: "svg",
				},
			},
			want: "testdata/tgae_golden.svg",
		},
	} {
		t.Run(tc.want, func(t *testing.T) {
			var resp PlotResponse
			testPlotS2(t, ts, tc.req, &resp)

			raw, err := base64.StdEncoding.DecodeString(resp.Data)
			if err != nil {
				t.Fatal(err)
			}

			if *cmpimg.GenerateTestData {
				ioutil.WriteFile(tc.want, raw, 0644)
			}

			want, err := ioutil.ReadFile(tc.want)
			if err != nil {
				t.Fatal(err)
			}

			typ := tc.req.Options.Type
			if typ == "" {
				typ = "png"
			}
			if ok, err := cmpimg.Equal(typ, raw, want); !ok || err != nil {
				ioutil.WriteFile(strings.Replace(tc.want, "_golden", "", -1), raw, 0644)
				t.Fatalf("reference files differ: err=%v ok=%v", err, ok)
			}
		})
	}
}

func testPlotS2(t *testing.T, ts *httptest.Server, req PlotS2Request, resp *PlotResponse) {
	t.Helper()

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(req)
	if err != nil {
		t.Fatalf("could not encode request: %v", err)
	}

	hreq, err := http.NewRequest(http.MethodPost, ts.URL+"/plot-s2", body)
	if err != nil {
		t.Fatalf("could not create http request: %v", err)
	}
	srv.addCookies(hreq)

	hresp, err := ts.Client().Do(hreq)
	if err != nil {
		t.Fatalf("could not post http request: %v", err)
	}
	defer hresp.Body.Close()

	if hresp.StatusCode != http.StatusOK {
		t.Fatalf("could not plot h1: %v", hresp.StatusCode)
	}

	err = json.NewDecoder(hresp.Body).Decode(resp)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
}

func TestPlotTree(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	const uri = "https://github.com/go-hep/hep/raw/master/groot/testdata/small-flat-tree.root"
	testOpenFile(t, ts, uri, http.StatusOK)
	defer testCloseFile(t, ts, uri)

	for _, tc := range []struct {
		req  PlotTreeRequest
		want string
	}{
		{
			req: PlotTreeRequest{
				URI:  uri,
				Obj:  "tree",
				Vars: []string{"Int32"},
			},
			want: "testdata/tree_i32_golden.png",
		},
		{
			req: PlotTreeRequest{
				URI:  uri,
				Dir:  "/",
				Obj:  "tree",
				Vars: []string{"Float64"},
			},
			want: "testdata/tree_f64_golden.png",
		},
		{
			req: PlotTreeRequest{
				URI:  uri,
				Dir:  "/",
				Obj:  "tree",
				Vars: []string{"ArrayFloat64"},
			},
			want: "testdata/tree_array_f64_golden.png",
		},
		{
			req: PlotTreeRequest{
				URI:  uri,
				Dir:  "/",
				Obj:  "tree",
				Vars: []string{"SliceFloat64"},
			},
			want: "testdata/tree_slice_f64_golden.png",
		},
	} {
		t.Run(tc.want, func(t *testing.T) {
			var resp PlotResponse
			testPlotTree(t, ts, tc.req, &resp)

			raw, err := base64.StdEncoding.DecodeString(resp.Data)
			if err != nil {
				t.Fatal(err)
			}

			if *cmpimg.GenerateTestData {
				ioutil.WriteFile(tc.want, raw, 0644)
			}

			want, err := ioutil.ReadFile(tc.want)
			if err != nil {
				t.Fatal(err)
			}

			typ := tc.req.Options.Type
			if typ == "" {
				typ = "png"
			}
			if ok, err := cmpimg.Equal(typ, raw, want); !ok || err != nil {
				ioutil.WriteFile(strings.Replace(tc.want, "_golden", "", -1), raw, 0644)
				t.Fatalf("reference files differ: err=%v ok=%v", err, ok)
			}
		})
	}
}

func testPlotTree(t *testing.T, ts *httptest.Server, req PlotTreeRequest, resp *PlotResponse) {
	t.Helper()

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(req)
	if err != nil {
		t.Fatalf("could not encode request: %v", err)
	}

	hreq, err := http.NewRequest(http.MethodPost, ts.URL+"/plot-tree", body)
	if err != nil {
		t.Fatalf("could not create http request: %v", err)
	}
	srv.addCookies(hreq)

	hresp, err := ts.Client().Do(hreq)
	if err != nil {
		t.Fatalf("could not post http request: %v", err)
	}
	defer hresp.Body.Close()

	if hresp.StatusCode != http.StatusOK {
		t.Fatalf("could not plot h1: %v", hresp.StatusCode)
	}

	err = json.NewDecoder(hresp.Body).Decode(resp)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
}

func (srv *Server) addCookies(req *http.Request) {
	for _, cookie := range srv.cookies {
		req.AddCookie(cookie)
	}
}

func setupCookie(srv *Server) {
	v, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	cookie := &http.Cookie{
		Name:    cookieName,
		Value:   v,
		Expires: time.Now().Add(24 * time.Hour),
	}
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.sessions[cookie.Value] = NewDB(filepath.Join(srv.dir, cookie.Value))
	srv.cookies[cookie.Value] = cookie
}
