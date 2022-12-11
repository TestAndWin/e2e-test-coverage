/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/TestAndWin/e2e-coverage/coverage/reporter"
	"github.com/gin-gonic/gin"
)

// UploadMochaSummaryReport godoc
// @Summary      Adds test result of a mocha summary report
// @Description  Adds test result of a mocha summary report.
// @Tags         mocha
// @Produce      json
// @Param        id    path      int     true  "Product ID"
// @Param        test  body      string  true  "Mocha JSON"
// @Success      201   object string
// @Router       /coverage/:id/upload-mocha-summary-report [POST]
func UploadMochaSummaryReport(c *gin.Context) {
	testResults, err := reporter.ReadMochaResultFromContext(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		var status []string
		for _, tr := range testResults {
			uploaded, err := repo.HasTestBeenUploaded(tr.Uuid)
			if err != nil {
				log.Println(err)
				status = append(status, err.Error())
			} else {
				if uploaded {
					status = append(status, tr.Uuid+" already uploaded")
				} else {
					// Check if product, area and feature exist
					pid := c.Param("id")
					aid, fid, err := repo.GetAreaAndFeatureId(tr.Area, tr.Feature, pid)
					if err != nil && err != sql.ErrNoRows {
						log.Println(err)
						status = append(status, err.Error())
					} else {
						// TODO Url
						var id int64
						var err error
						if aid != 0 && fid != 0 {
							id, err = repo.InsertTestResult(pid, aid, fid, tr)
						} else {
							id, err = repo.InsertTestResultWithoutAreaFeature(pid, tr)
						}
						if err != nil {
							status = append(status, err.Error())
						} else {
							status = append(status, strconv.FormatInt(id, 10))
						}
					}
				}
			}
		}
		c.JSON(http.StatusCreated, status)
	}
}
