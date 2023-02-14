// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"path"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/xrdfs"
)

func testFile_Close(t *testing.T, addr string) {
	t.Parallel()

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	file, err := fs.Open(context.Background(), "/tmp/dir1/file1.txt", xrdfs.OpenModeOtherRead, xrdfs.OpenOptionsNone)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}

	err = file.Close(context.Background())
	if err != nil {
		t.Fatalf("invalid close call: %v", err)
	}
}

func TestFile_Close(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFile_Close(t, addr)
		})
	}
}

func testFile_CloseVerify(t *testing.T, addr string) {
	t.Parallel()

	fileName := "close-verify.txt"
	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	dir, err := tempdir(client, "/tmp/", "xrd-test-close-verify")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = fs.RemoveAll(context.Background(), dir)
	}()
	filePath := path.Join(dir, fileName)

	file, err := fs.Open(context.Background(), filePath, xrdfs.OpenModeOwnerWrite, xrdfs.OpenOptionsNew)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}

	// TODO: Remove these 2 lines when XRootD server will follow protocol specification and fail such requests.
	// See https://github.com/xrootd/xrootd/issues/727.
	defer file.Close(context.Background())
	t.Skip("Skipping test because the XRootD C++ server doesn't fail request when the wrong size is passed.")

	err = file.CloseVerify(context.Background(), 14)
	if err == nil {
		t.Fatal("close call should fail when the wrong size is passed")
	}
}

func TestFile_CloseVerify(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFile_CloseVerify(t, addr)
		})
	}
}

func testFile_ReadAt(t *testing.T, addr string) {
	t.Parallel()

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	file, err := fs.Open(context.Background(), "/tmp/file1.txt", xrdfs.OpenModeOtherRead, xrdfs.OpenOptionsNone)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer file.Close(context.Background())

	want := []byte("Hello XRootD.\n")
	got := make([]uint8, 20)
	n, err := file.ReadAt(got, 0)
	if err != nil {
		t.Fatalf("invalid read call: %v", err)
	}

	if n != len(want) {
		t.Fatalf("read count does not match:\ngot = %v\nwant = %v", n, len(want))
	}

	if !reflect.DeepEqual(got[:n], want) {
		t.Fatalf("read data does not match:\ngot = %v\nwant = %v", got[:n], want)
	}
}

func TestFile_ReadAt(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFile_ReadAt(t, addr)
		})
	}
}

func testFile_WriteAt(t *testing.T, addr string) {
	t.Parallel()

	fileName := "test_rw.txt"
	want := make([]byte, 8*1024)
	rand.Read(want)

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()
	fs := client.FS()

	dir, err := tempdir(client, "/tmp/", "xrd-test-file-write-at-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = fs.RemoveAll(context.Background(), dir)
	}()
	filePath := path.Join(dir, fileName)

	file, err := fs.Open(context.Background(), filePath, xrdfs.OpenModeOwnerWrite, xrdfs.OpenOptionsNew)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer file.Close(context.Background())

	_, err = file.WriteAt(want, 0)
	if err != nil {
		t.Fatalf("invalid write call: %v", err)
	}

	err = file.Sync(context.Background())
	if err != nil {
		t.Fatalf("invalid sync call: %v", err)
	}

	file.Close(context.Background())
	file, err = fs.Open(context.Background(), filePath, xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsOpenRead)
	if err != nil {
		t.Fatalf("could not open %q: %v", filePath, err)
	}
	defer file.Close(context.Background())

	got := make([]uint8, len(want)+10)
	n, err := file.ReadAt(got, 0)
	if err != nil {
		t.Fatalf("invalid read call: %v", err)
	}

	if n != len(want) {
		t.Fatalf("read count does not match:\ngot = %v\nwant = %v", n, len(want))
	}

	if !reflect.DeepEqual(got[:n], want) {
		t.Fatalf("read data does not match:\ngot = %v\nwant = %v", got[:n], want)
	}

}

func TestFile_WriteAt(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFile_WriteAt(t, addr)
		})
	}
}

func testFile_Truncate(t *testing.T, addr string) {
	t.Parallel()

	fileName := "test_truncate.txt"
	write := []uint8{1, 2, 3, 4, 5, 6, 7, 8}
	want := write[:4]

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	dir, err := tempdir(client, "/tmp/", "xrd-test-truncate")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = fs.RemoveAll(context.Background(), dir)
	}()
	filePath := path.Join(dir, fileName)

	file, err := fs.Open(context.Background(), filePath, xrdfs.OpenModeOwnerWrite, xrdfs.OpenOptionsNew)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer file.Close(context.Background())

	_, err = file.WriteAt(write, 0)
	if err != nil {
		t.Fatalf("invalid write call: %v", err)
	}

	err = file.Truncate(context.Background(), int64(len(want)))
	if err != nil {
		t.Fatalf("invalid truncate call: %v", err)
	}

	err = file.Sync(context.Background())
	if err != nil {
		t.Fatalf("invalid sync call: %v", err)
	}

	err = file.Close(context.Background())
	if err != nil {
		t.Fatalf("invalid close call: %v", err)
	}
	file, err = fs.Open(context.Background(), filePath, xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsOpenRead)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer file.Close(context.Background())

	got := make([]uint8, len(want)+10)
	n, err := file.ReadAt(got, 0)
	if err != nil {
		t.Fatalf("invalid read call: %v", err)
	}

	if n != len(want) {
		t.Fatalf("read count does not match:\ngot = %v\nwant = %v", n, len(want))
	}

	if !reflect.DeepEqual(got[:n], want) {
		t.Fatalf("read data does not match:\ngot = %v\nwant = %v", got[:n], want)
	}

}

func TestFile_Truncate(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFile_Truncate(t, addr)
		})
	}
}

func testFile_Stat(t *testing.T, addr string) {
	t.Parallel()

	want := fstest["/tmp/dir1/file1.txt"]

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()
	file, err := fs.Open(context.Background(), "/tmp/dir1/file1.txt", xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsNone)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer file.Close(context.Background())

	got, err := file.Stat(context.Background())
	if err != nil {
		t.Fatalf("invalid stat call: %v", err)
	}

	if !reflect.DeepEqual(&got, want) {
		t.Fatalf("stat info does not match:\ngot = %v\nwant = %v", &got, want)
	}
	if !reflect.DeepEqual(file.Info(), want) {
		t.Fatalf("stat info does not match:\nfile.Info = %v\nwant = %v", file.Info(), want)
	}
}

func TestFile_Stat(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFile_Stat(t, addr)
		})
	}
}

func testFile_StatVirtualFS(t *testing.T, addr string) {
	t.Parallel()

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()
	file, err := fs.Open(context.Background(), "/tmp/dir1/file1.txt", xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsNone)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer file.Close(context.Background())

	// FIXME: Investigate whether this request is allowed by the protocol: https://github.com/xrootd/xrootd/issues/728
	t.Skip("Skipping this test because XRootD server probably doesn't support such requests.")

	want := xrdfs.VirtualFSStat{
		NumberRW:      1,
		FreeRW:        444,
		UtilizationRW: 6,
	}

	got, err := file.StatVirtualFS(context.Background())
	if err != nil {
		t.Fatalf("invalid stat call: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("File.FetchVirtualStatInfo()\ngot = %v\nwant = %v", got, want)
	}
}

func TestFile_StatVirtualFS(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFile_StatVirtualFS(t, addr)
		})
	}
}

func testFile_VerifyWriteAt(t *testing.T, addr string) {
	t.Parallel()

	// TODO: Enable this test once XRootD server starts to support kXR_verifyw request: https://github.com/xrootd/xrootd/issues/738.
	t.Skipf("Skipping this test because XRootD server doesn't support such request.")

	fileName := "test_verify_write.txt"
	want := make([]byte, 8*1024)
	rand.Read(want)

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()
	fs := client.FS()

	dir, err := tempdir(client, "/tmp/", "xrd-test-verify-write")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = fs.RemoveAll(context.Background(), dir)
	}()
	filePath := path.Join(dir, fileName)

	file, err := fs.Open(context.Background(), filePath, xrdfs.OpenModeOwnerWrite, xrdfs.OpenOptionsNew)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer file.Close(context.Background())

	err = file.VerifyWriteAt(context.Background(), want, 0)
	if err != nil {
		t.Fatalf("invalid verifyw call: %v", err)
	}

	err = file.Sync(context.Background())
	if err != nil {
		t.Fatalf("invalid sync call: %v", err)
	}

	err = file.Close(context.Background())
	if err != nil {
		t.Fatalf("invalid close call: %v", err)
	}
	file, err = fs.Open(context.Background(), filePath, xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsOpenRead)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer file.Close(context.Background())

	got := make([]uint8, len(want)+10)
	n, err := file.ReadAt(got, 0)
	if err != nil {
		t.Fatalf("invalid read call: %v", err)
	}

	if n != len(want) {
		t.Fatalf("read count does not match:\ngot = %v\nwant = %v", n, len(want))
	}

	if !reflect.DeepEqual(got[:n], want) {
		t.Fatalf("read data does not match:\ngot = %v\nwant = %v", got[:n], want)
	}

}

func TestFile_VerifyWriteAt(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFile_VerifyWriteAt(t, addr)
		})
	}
}

func ExampleClient_write() {
	ctx := context.Background()
	const username = "gopher"
	client, err := NewClient(ctx, "ccxrootdgotest.in2p3.fr:9001", username)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	file, err := client.FS().Open(ctx, "/tmp/test.txt", xrdfs.OpenModeOwnerRead|xrdfs.OpenModeOwnerWrite, xrdfs.OpenOptionsOpenUpdate|xrdfs.OpenOptionsNew)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close(ctx)

	if _, err := file.WriteAt([]byte("test"), 0); err != nil {
		log.Fatal(err)
	}

	if err := file.Sync(ctx); err != nil {
		log.Fatal(err)
	}

	if err := file.Close(ctx); err != nil {
		log.Fatal(err)
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_read() {
	ctx := context.Background()
	const username = "gopher"
	client, err := NewClient(ctx, "ccxrootdgotest.in2p3.fr:9001", username)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	file, err := client.FS().Open(ctx, "/tmp/test.txt", xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsOpenRead)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close(ctx)

	data := make([]byte, 10)
	n, err := file.ReadAt(data, 0)
	if err != nil {
		log.Fatal(err)
	}

	data = data[:n]
	fmt.Printf("%s\n", data)

	if err := file.Close(ctx); err != nil {
		log.Fatal(err)
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}

	// Output:
	// test
}
