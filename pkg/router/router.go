/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package router

import (
	"net/http"
	"os"

	"github.com/TestAndWin/e2e-coverage/pkg/controller"
	"github.com/TestAndWin/e2e-coverage/ui"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func HandleIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		html := ui.MustAsset("index.html")
		c.Data(200, "text/html; charset=UTF-8", html)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func HandleRequest() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(CORSMiddleware())

	// VUE routes
	h := gin.WrapH(http.FileServer(ui.AssetFile()))
	router.GET("/favicon.ico", h)
	router.GET("/js/*filepath", h)
	router.GET("/css/*filepath", h)
	router.GET("/img/*filepath", h)
	router.GET("/fonts/*filepath", h)
	router.NoRoute(HandleIndex())

	// Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API
	v1 := router.Group("/api/v1")
	{
		v1.POST("/products", controller.AddProduct)
		v1.GET("/products", controller.GetProducts)
		v1.PUT("/products/:id", controller.UpdateProduct)
		v1.DELETE("/products/:id", controller.DeleteProduct)

		v1.POST("/areas", controller.AddArea)
		v1.GET("/products/:id/areas", controller.GetProductAreas)
		v1.PUT("/areas/:id", controller.UpdateArea)
		v1.DELETE("/areas/:id", controller.DeleteArea)

		v1.POST("/features", controller.AddFeature)
		v1.GET("/areas/:id/features", controller.GetAreaFeatures)
		v1.PUT("/features/:id", controller.UpdateFeature)
		v1.DELETE("/features/:id", controller.DeleteFeature)

		v1.POST("/tests", controller.AddTest)
		v1.GET("/tests", controller.GetAllTestForSuiteFile)
		v1.DELETE("/tests/:id", controller.DeleteTest)

		v1.POST("/expl-tests", controller.AddExplTest)
		v1.GET("/expl-tests/area/:areaid", controller.GetExplTestsForArea)
		v1.DELETE("/expl-tests/:id", controller.DeleteExplTest)

		v1.POST("/coverage/:id/upload-mocha-report", controller.UploadMochaReport)
		v1.GET("/coverage/:id/areas", controller.GetAreaCoverage)
		v1.GET("/coverage/areas/:id/features", controller.GetFeatureCoverage)
		v1.GET("/coverage/features/:id/tests", controller.GetTestsCoverage)
	}

	port, found := os.LookupEnv("PORT")
	if found {
		router.Run("localhost:" + port)
	} else {
		router.Run("localhost:8080")
	}
}
