/*
Copyright (c) 2022-2024, webmaster@testandwin.net, Michael Schlottmann
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
	Title  string  `json:"title"`
	Tests  []Test  `json:"tests"`
	Suites []Suite `json:"suites"`
}

type Test struct {
	Pass    bool `json:"pass"`
	Fail    bool `json:"fail"`
	Pending bool `json:"pending"`
	Skipped bool `json:"skipped"`
}

// Iterate through the mocha report and get all the needed data. Currently it is only supported, that the results section contains only one suite entry.
func ReadMochaResultFromContext(c *gin.Context) ([]TestResult, error) {
	var m Mocha
	err := c.BindJSON(&m)
	if err != nil {
		return nil, fmt.Errorf("error during BindJSON(): %w", err)
	}
	return getTestResultFromMocha(m)
}

func getTestResultFromMocha(m Mocha) ([]TestResult, error) {
	parsedEndTime, err := time.Parse("2006-01-02T15:04:05.000Z", m.Stats.End)
	if err != nil {
		return nil, fmt.Errorf("error parsing end time: %w", err)
	}

	var results []TestResult
	for _, result := range m.Results {
		tr := extractSuiteDetails(result.Suites)
		tr.File = result.File
		tr.Uuid = result.Uuid
		tr.TestRun = parsedEndTime

		total, passes, pending, failures, skipped := countTests(result.Suites)
		tr.Total = total
		tr.Passes = passes
		tr.Pending = pending
		tr.Failures = failures
		tr.Skipped = skipped
		results = append(results, tr)
	}
	return results, nil
}

func extractSuiteDetails(suites []Suite) TestResult {
	tr := TestResult{}
	if len(suites) > 0 {
		if parts := strings.Split(suites[0].Title, "|"); len(parts) > 2 {
			tr.Area = parts[0]
			tr.Feature = parts[1]
			tr.Suite = parts[2]
		} else {
			tr.Suite = suites[0].Title
		}
	}
	return tr
}

func countTests(suites []Suite) (total, passes, pending, failures, skipped int) {
	for _, suite := range suites {
		tests := appendSuiteTests(suite)

		for _, test := range tests {
			total++
			switch {
			case test.Pass:
				passes++
			case test.Fail:
				failures++
			case test.Skipped:
				skipped++
			case test.Pending:
				pending++
			}
		}
	}
	return
}

func appendSuiteTests(suite Suite) []Test {
	tests := suite.Tests
	for _, subSuite := range suite.Suites {
		tests = append(tests, subSuite.Tests...)
	}
	return tests
}
