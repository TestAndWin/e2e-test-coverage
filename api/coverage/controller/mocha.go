/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/TestAndWin/e2e-coverage/coverage/reporter"
	"github.com/TestAndWin/e2e-coverage/errors"
	"github.com/TestAndWin/e2e-coverage/response"
	"github.com/gin-gonic/gin"
)

// UploadMochaSummaryReport godoc
// @Summary      Add test results of a mocha summary report
// @Description  Add test results of a mocha summary report.
// @Tags         mocha
// @Produce      json
// @Param        id            path      int     true   "Product ID"
// @Param        apiKey        header    string  true   "Api Key"
// @Param        testReportUrl header    string  false  "Url of the detail test report"
// @Param        test          body      string  true   "Mocha JSON"
// @Success      201  {object} string
// @Failure      400  {string}  ErrorResponse
// @Router       /coverage/:id/upload-mocha-summary-report [POST]
func UploadMochaSummaryReport(c *gin.Context) {
	testResults, err := reporter.ReadMochaResultFromContext(c)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Error reading Mocha result", err))
		return
	}

	pid := c.Param("id")
	testReportUrl := c.GetHeader("testReportUrl")
	component := c.GetHeader("component")

	status, err := processTestResults(testResults, pid, testReportUrl, component)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}

	response.Created(c, status)
}

func processTestResults(testResults []reporter.TestResult, pid, testReportUrl, component string) ([]string, error) {
	var status []string
	for _, tr := range testResults {
		resultStatus, err := processTestResult(tr, pid, testReportUrl, component)
		if err != nil {
			log.Printf("Error processing test result: %v", err)
			status = append(status, err.Error())
		} else {
			status = append(status, resultStatus)
		}
	}
	return status, nil
}

func processTestResult(tr reporter.TestResult, pid, testReportUrl, component string) (string, error) {
	repo := getRepository()
	
	uploaded, err := repo.HasTestBeenUploaded(tr.Uuid)
	if err != nil {
		return "", fmt.Errorf("error checking if test was uploaded: %w", err)
	}
	if uploaded {
		return tr.Uuid + " already uploaded", nil
	}

	aid, fid, err := repo.GetAreaAndFeatureId(tr.Area, tr.Feature, pid)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("error getting area and feature ID: %w", err)
	}

	isFirst, err := repo.IsThisTheFirstUpload(pid, aid, fid, tr.Suite, tr.File, component)
	if err != nil {
		return "", fmt.Errorf("error checking if this is the first upload: %w", err)
	}

	var id int64
	if aid != 0 && fid != 0 {
		id, err = repo.InsertTestResult(pid, aid, fid, component, testReportUrl, isFirst, tr)
	} else {
		id, err = repo.InsertTestResultWithoutAreaFeature(pid, component, testReportUrl, isFirst, tr)
	}
	if err != nil {
		return "", fmt.Errorf("error inserting test result: %w", err)
	}

	return strconv.FormatInt(id, 10), nil
}
