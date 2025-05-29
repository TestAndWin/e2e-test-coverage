/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"fmt"

	"github.com/TestAndWin/e2e-coverage/coverage/model"
	"github.com/TestAndWin/e2e-coverage/errors"
	"github.com/TestAndWin/e2e-coverage/response"
	"github.com/gin-gonic/gin"
)

// GetAreaCoverage godoc
// @Summary		   Get coverage for all product areas.
// @Description  Get coverage for all product areas. Only tests from the last 28 days are considered.
// @Tags         coverage
// @Produce      json
// @Param        product    path      int     true  "Product ID"
// @Success      200  {array}  model.Area
// @Failure      400  {string}  ErrorResponse
// @Router       /api/v1/coverage/{id}/areas [GET]
func GetAreaCoverage(c *gin.Context) {
	pId := c.Param("id")
	repo, err := getRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	areas, err := repo.GetAllProductAreas(pId)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Unable to get all product areas", err))
		return
	}

	tests, err := repo.GetAreaCoverageForProduct(pId)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Unable to get area coverage for product", err))
		return
	}

	areasCoverage, err := processAreaCoverage(areas, tests)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Unable to process area coverage", err))
		return
	}
	// Return consistent response format
	if len(areasCoverage) == 0 {
		response.EmptyList(c)
	} else {
		response.ResponseWithDataAndCount(c, 200, areasCoverage, int64(len(areasCoverage)))
	}
}

func processAreaCoverage(areas []model.Area, tests map[int64]model.Test) ([]model.Area, error) {
	repo, err := getRepository()
	if err != nil {
		return nil, err
	}
	areasCoverage := []model.Area{}
	for _, a := range areas {
		// Iterate through all areas and check if there is coverage data for that area
		if t, ok := tests[a.Id]; ok {
			a.Total = t.Total
			a.Passes = t.Passes
			a.Pending = t.Pending
			a.Failures = t.Failures
			a.Skipped = t.Skipped
			a.FirstTotal = t.FirstTotal
		}
		// Add expl. tests
		explTest, err := repo.GetExplTestOverviewForArea(a.Id)
		if err != nil {
			return nil, fmt.Errorf("error fetching exploratory tests for area %d: %w", a.Id, err)
		} else {
			a.ExplTests = explTest.ExplTests
			a.ExplRating = explTest.ExplRating
		}
		areasCoverage = append(areasCoverage, a)
	}
	return areasCoverage, nil
}

// GetFeatureCoverage godoc
// @Summary      Get coverage for all area features.
// @Description  Get coverage for all area features. Only tests from the last 28 days are considered.
// @Tags         coverage
// @Produce      json
// @Param        product    path      int     true  "Area ID"
// @Success      200  {array}  model.Feature
// @Failure      400  {string}  ErrorResponse
// @Router       /api/v1/coverage/areas/{id}/features [get]
func GetFeatureCoverage(c *gin.Context) {
	repo, err := getRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	features, err := repo.GetAllAreaFeatures(c.Param("id"))
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Error getting area features", err))
		return
	}

	tests, err := repo.GetFeatureCoverageForArea(c.Param("id"))
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Error getting feature coverage", err))
		return
	}

	featuresCoverage := []model.Feature{}
	for _, f := range features {
		if t, ok := tests[f.Id]; ok {
			f.Total = t.Total
			f.Passes = t.Passes
			f.Pending = t.Pending
			f.Failures = t.Failures
			f.Skipped = t.Skipped
			f.FirstTotal = t.FirstTotal
		}
		featuresCoverage = append(featuresCoverage, f)
	}
	response.OK(c, featuresCoverage)
}

// GetTestsCoverage godoc
// @Summary      Get coverage for all tests of a feature.
// @Description  Get coverage for all tests of a feature for the last 28 days.
// @Tags         coverage
// @Produce      json
// @Param        id    path      int     true  "Feature ID"
// @Success      200  {array}  model.Test
// @Failure      400  {string}  ErrorResponse
// @Router       /coverage/features/:id/tests [get]
func GetTestsCoverage(c *gin.Context) {
	repo, err := getRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	t, err := repo.GetAllFeatureTests(c.Param("id"))
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Error getting feature tests", err))
		return
	}
	response.OK(c, t)
}

// GetProductTestsCoverage godoc
// @Summary      Get coverage for all tests of a product.
// @Description  Get coverage for all tests of a product for the last 28 days.
// @Tags         coverage
// @Produce      json
// @Param        id    path      int     true  "Product ID"
// @Success      200  {array}  model.Test
// @Failure      400  {string}  ErrorResponse
// @Router       /coverage/products/:id/tests [get]
func GetProductTestsCoverage(c *gin.Context) {
	repo, err := getRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	t, err := repo.GetAllProductTests(c.Param("id"))
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Error getting product tests", err))
		return
	}
	response.OK(c, t)
}

// GetComponents godoc
// @Summary      Get all components with their latest test run
// @Description  Get all components with their latest test run
// @Tags         coverage
// @Produce      json
// @Success      200  {array}  model.Test
// @Failure      400  {string}  ErrorResponse
// @Router       /coverage/components [get]
func GetComponents(c *gin.Context) {
	repo, err := getRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	t, err := repo.GetComponents()
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Error getting components", err))
		return
	}
	response.OK(c, t)
}
