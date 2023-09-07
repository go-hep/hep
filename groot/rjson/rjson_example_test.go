// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/rjson"
	"go-hep.org/x/hep/hbook"
)

func ExampleMarshal() {
	h := hbook.NewH1D(5, 0, 5)
	h.Fill(1, 1)
	h.Fill(-1, 1)
	h.Ann["name"] = "h1"
	h.Ann["title"] = "my title"

	raw, err := rjson.Marshal(rhist.NewH1FFrom(h))
	if err != nil {
		log.Fatalf("could not marshal to ROOT JSON: %+v", err)
	}

	buf := new(bytes.Buffer)
	err = json.Indent(buf, raw, "", "  ")
	if err != nil {
		log.Fatalf("could not indent JSON: %+v", err)
	}

	fmt.Printf("json: %s\n", buf.String())

	// Output:
	// json: {
	//   "_typename": "TH1F",
	//   "fUniqueID": 0,
	//   "fBits": 50331648,
	//   "fName": "h1",
	//   "fTitle": "my title",
	//   "fLineColor": 602,
	//   "fLineStyle": 1,
	//   "fLineWidth": 1,
	//   "fFillColor": 0,
	//   "fFillStyle": 1001,
	//   "fMarkerColor": 1,
	//   "fMarkerStyle": 1,
	//   "fMarkerSize": 1,
	//   "fNcells": 7,
	//   "fXaxis": {
	//     "_typename": "TAxis",
	//     "fUniqueID": 0,
	//     "fBits": 50331648,
	//     "fName": "xaxis",
	//     "fTitle": "",
	//     "fNdivisions": 510,
	//     "fAxisColor": 1,
	//     "fLabelColor": 1,
	//     "fLabelFont": 42,
	//     "fLabelOffset": 0.005,
	//     "fLabelSize": 0.035,
	//     "fTickLength": 0.03,
	//     "fTitleOffset": 1,
	//     "fTitleSize": 0.035,
	//     "fTitleColor": 1,
	//     "fTitleFont": 42,
	//     "fNbins": 5,
	//     "fXmin": 0,
	//     "fXmax": 5,
	//     "fXbins": [
	//       0,
	//       1,
	//       2,
	//       3,
	//       4,
	//       5
	//     ],
	//     "fFirst": 0,
	//     "fLast": 0,
	//     "fBits2": 0,
	//     "fTimeDisplay": false,
	//     "fTimeFormat": "",
	//     "fLabels": null,
	//     "fModLabs": null
	//   },
	//   "fYaxis": {
	//     "_typename": "TAxis",
	//     "fUniqueID": 0,
	//     "fBits": 50331648,
	//     "fName": "yaxis",
	//     "fTitle": "",
	//     "fNdivisions": 510,
	//     "fAxisColor": 1,
	//     "fLabelColor": 1,
	//     "fLabelFont": 42,
	//     "fLabelOffset": 0.005,
	//     "fLabelSize": 0.035,
	//     "fTickLength": 0.03,
	//     "fTitleOffset": 1,
	//     "fTitleSize": 0.035,
	//     "fTitleColor": 1,
	//     "fTitleFont": 42,
	//     "fNbins": 1,
	//     "fXmin": 0,
	//     "fXmax": 1,
	//     "fXbins": [],
	//     "fFirst": 0,
	//     "fLast": 0,
	//     "fBits2": 0,
	//     "fTimeDisplay": false,
	//     "fTimeFormat": "",
	//     "fLabels": null,
	//     "fModLabs": null
	//   },
	//   "fZaxis": {
	//     "_typename": "TAxis",
	//     "fUniqueID": 0,
	//     "fBits": 50331648,
	//     "fName": "zaxis",
	//     "fTitle": "",
	//     "fNdivisions": 510,
	//     "fAxisColor": 1,
	//     "fLabelColor": 1,
	//     "fLabelFont": 42,
	//     "fLabelOffset": 0.005,
	//     "fLabelSize": 0.035,
	//     "fTickLength": 0.03,
	//     "fTitleOffset": 1,
	//     "fTitleSize": 0.035,
	//     "fTitleColor": 1,
	//     "fTitleFont": 42,
	//     "fNbins": 1,
	//     "fXmin": 0,
	//     "fXmax": 1,
	//     "fXbins": [],
	//     "fFirst": 0,
	//     "fLast": 0,
	//     "fBits2": 0,
	//     "fTimeDisplay": false,
	//     "fTimeFormat": "",
	//     "fLabels": null,
	//     "fModLabs": null
	//   },
	//   "fBarOffset": 0,
	//   "fBarWidth": 1000,
	//   "fEntries": 2,
	//   "fTsumw": 2,
	//   "fTsumw2": 2,
	//   "fTsumwx": 0,
	//   "fTsumwx2": 2,
	//   "fMaximum": -1111,
	//   "fMinimum": -1111,
	//   "fNormFactor": 0,
	//   "fContour": [],
	//   "fSumw2": [
	//     1,
	//     0,
	//     1,
	//     0,
	//     0,
	//     0,
	//     0
	//   ],
	//   "fOption": "",
	//   "fFunctions": {
	//     "_typename": "TList",
	//     "name": "",
	//     "arr": [],
	//     "opt": []
	//   },
	//   "fBufferSize": 0,
	//   "fBuffer": [],
	//   "fBinStatErrOpt": 0,
	//   "fStatOverflows": 2,
	//   "fArray": [
	//     1,
	//     0,
	//     1,
	//     0,
	//     0,
	//     0,
	//     0
	//   ]
	// }
}
