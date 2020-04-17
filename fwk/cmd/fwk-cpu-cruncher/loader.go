// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"os"

	"go-hep.org/x/hep/fwk/job"
)

func loadConfig(fname string, app *job.Job) {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var data struct {
		Algs []struct {
			Name    string   `json:"name"`
			Inputs  []string `json:"inputs"`
			Outputs []string `json:"outputs"`
			CPU     []int64  `json:"runtimes"`
			Wall    []int64  `json:"runtimes_all"`
		} `json:"algorithms"`
	}

	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		panic(err)
	}

	for _, alg := range data.Algs {
		app.Create(job.C{
			Type: "main.CPUCruncher",
			Name: alg.Name,
			Props: job.P{
				"Inputs":  alg.Inputs,
				"Outputs": alg.Outputs,
				"CPU":     alg.CPU,
			},
		})
	}
}
