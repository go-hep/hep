// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package groot provides a pure-go read/write-access to ROOT files.
//
// A typical usage is as follows:
//
//   f, err := groot.Open("ntup.root")
//   if err != nil {
//       log.Fatal(err)
//   }
//   defer f.Close()
//
//   obj, err := f.Get("tree")
//   if err != nil {
//       log.Fatal(err)
//   }
//   tree := obj.(rtree.Tree)
//   fmt.Printf("entries= %v\n", tree.Entries())
//
// More complete examples on how to iterate over the content of a Tree can
// be found in the examples attached to groot.TreeScanner and groot.Scanner:
// https://godoc.org/go-hep.org/x/hep/groot/rtree#pkg-examples
//
// Another possibility is to look at:
// https://godoc.org/go-hep.org/x/hep/groot/cmd/root-ls,
// a command that inspects the content of ROOT files.
//
//
// File layout
//
// ROOT files are a suite of consecutive data records.
// Each data record consists of a header part, called a TKey, and a payload
// whose content, length and meaning are described by the header.
// The current ROOT file format encodes all data in big endian.
//
// ROOT files initially only supported 32b addressing.
// Large files support (>4Gb) was added later on by migrating to a 64b addressing.
//
// The on-disk binary layout of a ROOT file header looks like this:
//
//       Type        | Record Name | Description
//  =================+=============+===========================================
//     [4]byte       | "root"      | Root file identifier
//     int32         | fVersion    | File format version
//     int32         | fBEGIN      | Pointer to first data record
//     int32 [int64] | fEND        | Pointer to first free word at the EOF
//     int32 [int64] | fSeekFree   | Pointer to FREE data record
//     int32         | fNbytesFree | Number of bytes in FREE data record
//     int32         | nfree       | Number of free data records
//     int32         | fNbytesName | Number of bytes in TNamed at creation time
//     byte          | fUnits      | Number of bytes for file pointers
//     int32         | fCompress   | Compression level and algorithm
//     int32 [int64] | fSeekInfo   | Pointer to TStreamerInfo record
//     int32         | fNbytesInfo | Number of bytes in TStreamerInfo record
//     [18]byte      | fUUID       | Universal Unique ID
//  =================+=============+===========================================
//
// This is followed by a sequence of data records, starting at the fBEGIN
// offset from the beginning of the file.
//
// The on-disk binary layout of a data record is:
//
//        Type     | Member Name | Description
//  ===============+=============+===========================================
//   int32         | Nbytes      | Length of compressed object (in bytes)
//   int16         | Version     | TKey version identifier
//   int32         | ObjLen      | Length of uncompressed object
//   int32         | Datime      | Date and time when object was written to file
//   int16         | KeyLen      | Length of the key structure (in bytes)
//   int16         | Cycle       | Cycle of key
//   int32 [int64] | SeekKey     | Pointer to record itself (consistency check)
//   int32 [int64] | SeekPdir    | Pointer to directory header
//   byte          | lname       | Number of bytes in the class name
//   []byte        | ClassName   | Object Class Name
//   byte          | lname       | Number of bytes in the object name
//   []byte        | Name        | Name of the object
//   byte          | lTitle      | Number of bytes in the object title
//   []byte        | Title       | Title of the object
//   []byte        | DATA        | Data bytes associated to the object
//  ===============+=============+===========================================
//
// The high-level on-disk representation of a ROOT file is thus:
//
//  +===============+ -- 0
//  |               |
//  |  File Header  |
//  |               |
//  +===============+ -- fBEGIN offset
//  |               |
//  | Record Header | -->-+
//  |               |     |
//  +---------------+     |
//  |               |     |
//  |  Record Data  |     | Reference to next Record
//  |    Payload    |     |
//  |               |     |
//  +===============+ <---+
//  |               |
//  | Record Header | -->-+
//  |               |     |
//  +---------------+     |
//  |               |     |
//  |  Record Data  |     | Reference to next Record
//  |    Payload    |     |
//  |               |     |
//  +===============+ <---+
//  |               |
//         ...
//
//  |               |
//  +===============+ -- fSeekInfo
//  |               |
//  | Record Header | -->-+
//  |               |     |
//  +---------------+     |
//  |               |     |
//  |  Record Data  |     | Reference to next Record
//  |    Payload    |     |
//  |               |     |
//  +===============+ <---+ -- fEND offset
//
// Data records payloads and how to deserialize them are described by a TStreamerInfo.
// The list of all TStreamerInfos that are used to interpret the content of
// a ROOT file is stored at the end of that ROOT file, at offset fSeekInfo.
//
//
// Data records
//
// Data records' payloads may be compressed.
// Detecting whether a payload is compressed is usually done by comparing
// the object length (ObjLen) field of the record header with the length
// of the compressed object (Nbytes) field.
// If they differ after having subtracted the record header length, then
// the payload has been compressed.
//
// A record data payload is itself split into multiple chunks of maximum
// size 16*1024*1024 bytes.
// Each chunk consists of:
//
//  - the chunk header,
//  - the chunk compressed payload.
//
// The chunk header:
//
//  - 3 bytes to identify the compression algorithm and version,
//  - 3 bytes to identify the deflated buffer size,
//  - 3 bytes to identify the inflated buffer size.
//
// Streamer informations
//
// Streamers describe how a given type, for a given version of that type, is
// written on disk.
// In C++/ROOT, a streamer is represented as a TStreamerInfo class that can
// give metadata about the type it's describing (version, name).
// When reading a file, all the streamer infos are read back in memory, from
// disk, by reading the data record at offset fSeekInfo.
// A streamer info is actually a list of streamer elements, one for each field
// and, in C++, base class (in Go, this is emulated as an embedded field.)
//
// TODO: groot can not write trees yet.
package groot // import "go-hep.org/x/hep/groot"
