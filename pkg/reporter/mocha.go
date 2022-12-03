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
	Tests []Test `json:"tests"`
}

type Test struct {
	Pass    bool `json:"pass"`
	Fail    bool `json:"fail"`
	Pending bool `json:"pending"`
	Skipped bool `json:"skipped"`
}

// Iterate through the mocha report and get all the needed data. Currently it is only support, that the results section contains only one suite entry.
func ReadMochaResultFromContext(c *gin.Context) ([]TestResult, error) {
	var m Mocha
	if err := c.BindJSON(&m); err != nil {
		fmt.Println("Error during BindJSON(): ", err)
		return []TestResult{}, err
	}
	return getTestResultFromMocha(m), nil
}

func getTestResultFromMocha(m Mocha) []TestResult {
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", m.Stats.End)

	tests := []TestResult{}
	for _, result := range m.Results {
		tr := TestResult{}
		if strings.Contains(result.Suites[0].Title, "|") {
			tr.Area = strings.Split(result.Suites[0].Title, "|")[0]
			tr.Feature = strings.Split(result.Suites[0].Title, "|")[1]
			tr.Suite = strings.Split(result.Suites[0].Title, "|")[2]
		} else {
			tr.Suite = result.Suites[0].Title
		}
		tr.File = result.File
		tr.Uuid = result.Uuid
		tr.TestRun = t

		// We have to iterate through the suites/tests to get the numbers
		total := 0
		passes := 0
		pending := 0
		failures := 0
		skipped := 0
		for _, test := range result.Suites[0].Tests {
			total += 1
			if test.Pass {
				passes += 1
			}
			if test.Fail {
				failures += 1
			}
			if test.Skipped {
				skipped += 1
			}
			if test.Pending {
				pending += 1
			}
		}
		tr.Total = total
		tr.Passes = passes
		tr.Pending = pending
		tr.Failures = failures
		tr.Skipped = skipped
		tests = append(tests, tr)
	}
	return tests
}
