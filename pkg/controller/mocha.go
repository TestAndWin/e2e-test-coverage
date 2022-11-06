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

	"github.com/TestAndWin/e2e-coverage/pkg/model"
	"github.com/TestAndWin/e2e-coverage/pkg/reporter"
	"github.com/gin-gonic/gin"
)

// UploadMochaReport godoc
// @Summary      Adds test result of a mocha report
// @Description  Adds test result of a mocha report and returns the ID of the stored test.
// @Tags         mocha
// @Produce      json
// @Param        id    path      int     true  "Product ID"
// @Param        test  body      string  true  "Mocha JSON"
// @Success      201   object string
// @Router       /coverage/:id/upload-mocha-report [POST]
func UploadMochaReport(c *gin.Context) {
	tr, err := reporter.ReadMochaResultFromContext(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		uploaded, err := repo.HasTestBeenUploaded(tr.Uuid)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		} else {
			if uploaded {
				c.JSON(http.StatusBadRequest, gin.H{"status": "Test had already been uploaded"})
			} else {
				// Check if product, area and feature exist
				pid := c.Param("id")
				aid, fid, err := repo.GetAreaAndFeatureId(tr.Area, tr.Feature, pid)
				if err != nil && err != sql.ErrNoRows {
					log.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
				} else {
					// TODO Url
					var id int64
					var err error
					if aid != 0 && fid != 0 {
						id, err = repo.ExecuteSql(model.INSERT_TEST, pid, aid, fid, tr.Suite, tr.File, "", tr.Total, tr.Passes, tr.Pending, tr.Failures, tr.Skipped, tr.Uuid, tr.TestRun)
					} else {
						id, err = repo.ExecuteSql(model.INSERT_TEST_NO_AREA_FEATURE, pid, tr.Suite, tr.File, "", tr.Total, tr.Passes, tr.Pending, tr.Failures, tr.Skipped, tr.Uuid, tr.TestRun)
					}
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
					} else {
						c.JSON(http.StatusCreated, gin.H{"test-id": id})
					}
				}
			}
		}
	}
}
