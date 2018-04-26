// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package open // import "go-hep.org/x/hep/xrootd/requests/open"

const RequestID uint16 = 3010

type Response struct {
	FileHandle [4]byte
}

type Mode uint16

const (
	ModeOwnerRead    Mode = 0x100
	ModeOwnerWrite   Mode = 0x080
	ModeOwnerExecute Mode = 0x040

	ModeGroupRead    Mode = 0x020
	ModeGroupWrite   Mode = 0x010
	ModeGroupExecute Mode = 0x008

	ModeOtherRead    Mode = 0x004
	ModeOtherWrite   Mode = 0x002
	ModeOtherExecute Mode = 0x001
)

type Options uint16

const (
	// OptionsNone specifies that file is opened without specific options
	OptionsNone Options = 0
	// OptionsCompress specifies that file is opened even when compressed
	OptionsCompress Options = 1
	// OptionsDelete specifies that file is opened deleting any existing file
	OptionsDelete Options = 2
	// OptionsForce specifies that file is opened ignoring  file usage rules
	OptionsForce Options = 4
	// OptionsNew specifies that file is opened only if it does not already exist
	OptionsNew Options = 8
	// OptionsOpenRead specifies that file is opened only for reading
	OptionsOpenRead Options = 16
	// OptionsOpenUpdate specifies that file is opened only for reading and writing
	OptionsOpenUpdate Options = 32
	// OptionsAsyncOpen specifies that file is opened for asynchronous i/o
	OptionsAsync Options = 64
	// OptionsRefresh specifies that cached information on the file's location need to be updated
	OptionsRefresh Options = 128
	// OptionsMkPath specifies that directory path is created if it does not already exist
	OptionsMkPath Options = 256
	// OptionsOpenAppend specifies that file is opened only for appending
	OptionsOpenAppend Options = 512
	// OptionsReturnStatus specifies that file status information should be returned in the response
	OptionsReturnStatus Options = 1024
	// OptionsReplica specifies that file is opened for replica creation
	OptionsReplica Options = 2048
	// OptionsPOSC specifies that Persist On Successful Close (POSC) processing should be enabled
	OptionsPOSC Options = 4096
	// OptionsNoWait specifies that file is opened only if it does not cause a wait
	OptionsNoWait Options = 8192
	// OptionsSequentiallyIO specifies that file will be read or written sequentially
	OptionsSequentiallyIO Options = 16384
)

type Request struct {
	Mode       Mode
	Options    Options
	Reserved   [12]byte
	PathLength int32
	Path       []byte
}

func NewRequest(path string, mode Mode, options Options) Request {
	var pathBytes = make([]byte, len(path))
	copy(pathBytes, path)

	return Request{mode, options, [12]byte{}, int32(len(path)), pathBytes}
}
