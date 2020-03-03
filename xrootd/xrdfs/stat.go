// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrdfs // import "go-hep.org/x/hep/xrootd/xrdfs"

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
)

// StatFlags identifies the entry's attributes.
type StatFlags int32

const (
	// StatIsFile indicates that entry is a regular file if no other flag is specified.
	StatIsFile StatFlags = 0
	// StatIsExecutable indicates that entry is either an executable file or a searchable directory.
	StatIsExecutable StatFlags = 1
	// StatIsDir indicates that entry is a directory.
	StatIsDir StatFlags = 2
	// StatIsOther indicates that entry is neither a file nor a directory.
	StatIsOther StatFlags = 4
	// StatIsOffline indicates that the file is not online (i. e., on disk).
	StatIsOffline StatFlags = 8
	// StatIsReadable indicates that read access to that entry is allowed.
	StatIsReadable StatFlags = 16
	// StatIsWritable indicates that write access to that entry is allowed.
	StatIsWritable StatFlags = 32
	// StatIsPOSCPending indicates that the file was created with kXR_posc and has not yet been successfully closed.
	// kXR_posc is an option of open request indicating that the "Persist On Successful Close" processing is enabled and
	// the file will be persisted only when it has been explicitly closed.
	StatIsPOSCPending StatFlags = 64
)

// EntryStat holds the entry name and the entry stat information.
type EntryStat struct {
	EntryName   string    // EntryName is the name of entry.
	HasStatInfo bool      // HasStatInfo indicates if the following stat information is valid.
	ID          int64     // ID is the OS-dependent identifier assigned to this entry.
	EntrySize   int64     // EntrySize is the decimal size of the entry.
	Flags       StatFlags // Flags identifies the entry's attributes.
	Mtime       int64     // Mtime is the last modification time in Unix time units.
}

// EntryStatFrom creates an EntryStat that represents same information as the provided info.
func EntryStatFrom(info os.FileInfo) EntryStat {
	es := EntryStat{
		EntryName:   info.Name(),
		EntrySize:   info.Size(),
		Mtime:       info.ModTime().Unix(),
		HasStatInfo: true,
	}
	if info.IsDir() {
		es.Flags |= StatIsDir
	}
	if info.Mode()&0400 != 0 {
		es.Flags |= StatIsReadable
	}
	if info.Mode()&0200 != 0 {
		es.Flags |= StatIsWritable
	}
	return es
}

// Name implements os.FileInfo.
func (es EntryStat) Name() string {
	return es.EntryName
}

// Size implements os.FileInfo.
func (es EntryStat) Size() int64 {
	return es.EntrySize
}

// ModTime implements os.FileInfo.
func (es EntryStat) ModTime() time.Time {
	return time.Unix(es.Mtime, 0)
}

// Sys implements os.FileInfo.
func (es EntryStat) Sys() interface{} {
	return nil
}

// Mode implements os.FileInfo.
func (es EntryStat) Mode() os.FileMode {
	var mode os.FileMode
	if es.IsDir() {
		mode |= os.ModeDir
	}
	if es.IsWritable() {
		mode |= 0222
	}
	if es.IsReadable() {
		mode |= 0444
	}
	return mode
}

// IsExecutable indicates whether this entry is either an executable file or a searchable directory.
func (es EntryStat) IsExecutable() bool {
	return es.Flags&StatIsExecutable != 0
}

// IsDir indicates whether this entry is a directory.
func (es EntryStat) IsDir() bool {
	return es.Flags&StatIsDir != 0
}

// IsOther indicates whether this entry is neither a file nor a directory.
func (es EntryStat) IsOther() bool {
	return es.Flags&StatIsOther != 0
}

// IsOffline indicates whether this the file is not online (i. e., on disk).
func (es EntryStat) IsOffline() bool {
	return es.Flags&StatIsOffline != 0
}

// IsReadable indicates whether this read access to that entry is allowed.
func (es EntryStat) IsReadable() bool {
	return es.Flags&StatIsReadable != 0
}

// IsWritable indicates whether this write access to that entry is allowed.
func (es EntryStat) IsWritable() bool {
	return es.Flags&StatIsWritable != 0
}

// IsPOSCPending indicates whether this the file was created with kXR_posc and has not yet been successfully closed.
// kXR_posc is an option of open request indicating that the "Persist On Successful Close" processing is enabled and
// the file will be persisted only when it has been explicitly closed.
func (es EntryStat) IsPOSCPending() bool {
	return es.Flags&StatIsPOSCPending != 0
}

// MarshalXrd implements xrdproto.Marshaler.
func (o EntryStat) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	if !o.HasStatInfo {
		return nil
	}

	idStr := strconv.Itoa(int(o.ID))
	sizeStr := strconv.Itoa(int(o.EntrySize))
	flagsStr := strconv.Itoa(int(o.Flags))
	mtimeStr := strconv.Itoa(int(o.Mtime))

	wBuffer.WriteBytes([]byte(idStr + " " + sizeStr + " " + flagsStr + " " + mtimeStr))
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *EntryStat) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	var buf []byte
	for rBuffer.Len() != 0 {
		b := rBuffer.ReadU8()
		if b == '\x00' || b == '\n' {
			break
		}
		buf = append(buf, b)
	}

	stats := bytes.Split(buf, []byte{' '})
	if len(stats) < 4 {
		return fmt.Errorf("xrootd: statinfo \"%s\" doesn't have enough fields, expected format is: \"id size flags modtime\"", buf)
	}

	id, err := strconv.Atoi(string(stats[0]))
	if err != nil {
		return err
	}
	size, err := strconv.Atoi(string(stats[1]))
	if err != nil {
		return err
	}
	flags, err := strconv.Atoi(string(stats[2]))
	if err != nil {
		return err
	}
	mtime, err := strconv.Atoi(string(stats[3]))
	if err != nil {
		return err
	}

	o.HasStatInfo = true
	o.ID = int64(id)
	o.EntrySize = int64(size)
	o.Mtime = int64(mtime)
	o.Flags = StatFlags(flags)

	return nil
}

// VirtualFSStat holds the virtual file system information.
type VirtualFSStat struct {
	NumberRW           int // NumberRW is the number of nodes that can provide read/write space.
	FreeRW             int // FreeRW is the size, in megabytes, of the largest contiguous area of read/write free space.
	UtilizationRW      int // UtilizationRW is the percent utilization of the partition represented by FreeRW.
	NumberStaging      int // NumberStaging is the number of nodes that can provide staging space.
	FreeStaging        int // FreeStaging is the size, in megabytes, of the largest contiguous area of staging free space.
	UtilizationStaging int // UtilizationStaging is the percent utilization of the partition represented by FreeStaging.
}

// MarshalXrd implements xrdproto.Marshaler
func (o VirtualFSStat) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	nrw := strconv.Itoa(o.NumberRW)
	frw := strconv.Itoa(o.FreeRW)
	urw := strconv.Itoa(o.UtilizationRW)
	nstg := strconv.Itoa(o.NumberStaging)
	fstg := strconv.Itoa(o.FreeStaging)
	ustg := strconv.Itoa(o.UtilizationStaging)
	wBuffer.WriteBytes([]byte(nrw + " " + frw + " " + urw + " " + nstg + " " + fstg + " " + ustg))
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler
func (o *VirtualFSStat) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	var buf []byte
	for rBuffer.Len() != 0 {
		b := rBuffer.ReadU8()
		if b == '\x00' || b == '\n' {
			break
		}
		buf = append(buf, b)
	}

	stats := bytes.Split(buf, []byte{' '})
	if len(stats) < 6 {
		return fmt.Errorf("xrootd: virtual statinfo \"%s\" doesn't have enough fields, expected format is: \"nrw frw urw nstg fstg ustg\"", buf)
	}

	nrw, err := strconv.Atoi(string(stats[0]))
	if err != nil {
		return err
	}
	frw, err := strconv.Atoi(string(stats[1]))
	if err != nil {
		return err
	}
	urw, err := strconv.Atoi(string(stats[2]))
	if err != nil {
		return err
	}
	nstg, err := strconv.Atoi(string(stats[3]))
	if err != nil {
		return err
	}
	fstg, err := strconv.Atoi(string(stats[4]))
	if err != nil {
		return err
	}
	ustg, err := strconv.Atoi(string(stats[5]))
	if err != nil {
		return err
	}

	o.NumberRW = nrw
	o.FreeRW = frw
	o.UtilizationRW = urw
	o.NumberStaging = nstg
	o.FreeStaging = fstg
	o.UtilizationStaging = ustg

	return nil
}

var (
	_ os.FileInfo = (*EntryStat)(nil)
)
