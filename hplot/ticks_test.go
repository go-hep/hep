// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestTicks(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleTicks, t, "ticks.png")
}

func TestTicksTimeYearly(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleTicks_yearly, t, "timeseries_yearly.png")
}

func TestTicksTimeMonthly(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleTicks_monthly, t, "timeseries_monthly.png")
}

func TestTicksTimeDaily(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleTicks_daily, t, "timeseries_daily.png")
}
