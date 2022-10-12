/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package reporter

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Mocha struct {
	Stats   Stats     `json:"stats"`
	Results []Results `json:"results"`
}

type Stats struct {
	Tests    int    `json:"tests"`
	Passes   int    `json:"passes"`
	Pending  int    `json:"pending"`
	Failures int    `json:"failures"`
	Skipped  int    `json:"skipped"`
	End      string `json:"end"`
}

type Results struct {
	File   string  `json:"file"`
	Uuid   string  `json:"uuid"`
	Suites []Suite `json:"suites"`
}

type Suite struct {
	Title string `json:"title"`
	File  string `json:"file"`
}

func ReadMochaResultFromContext(c *gin.Context) (TestResult, error) {
	var m Mocha
	if err := c.BindJSON(&m); err != nil {
		fmt.Println("Error during BindJSON(): ", err)
		return TestResult{}, err
	}
	return getTestResultFromMocha(m), nil
}

func getTestResultFromMocha(m Mocha) TestResult {
	tr := TestResult{}

	if strings.Contains(m.Results[0].Suites[0].Title, "|") {
		tr.Area = strings.Split(m.Results[0].Suites[0].Title, "|")[0]
		tr.Feature = strings.Split(m.Results[0].Suites[0].Title, "|")[1]
		tr.Suite = strings.Split(m.Results[0].Suites[0].Title, "|")[2]
	} else {
		tr.Suite = m.Results[0].Suites[0].Title
	}

	if len(m.Results[0].Suites[0].File) > 0 {
		tr.File = m.Results[0].Suites[0].File
	} else {
		tr.File = m.Results[0].File
	}
	tr.Total = m.Stats.Tests
	tr.Passes = m.Stats.Passes
	tr.Pending = m.Stats.Pending
	tr.Failures = m.Stats.Failures
	tr.Skipped = m.Stats.Skipped
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", m.Stats.End)
	tr.TestRun = t
	tr.Uuid = m.Results[0].Uuid

	return tr
}
