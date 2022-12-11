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

	"github.com/TestAndWin/e2e-coverage/coverage/controller"
	"github.com/TestAndWin/e2e-coverage/ui"
	usercontroller "github.com/TestAndWin/e2e-coverage/user/controller"
	"github.com/TestAndWin/e2e-coverage/user/model"
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
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func AuthApi() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		v1.POST("/products", usercontroller.AuthUser(model.EDITOR), controller.AddProduct)
		v1.GET("/products", usercontroller.AuthUser(model.EDITOR), controller.GetProducts)
		v1.PUT("/products/:id", usercontroller.AuthUser(model.EDITOR), controller.UpdateProduct)
		v1.DELETE("/products/:id", usercontroller.AuthUser(model.EDITOR), controller.DeleteProduct)

		v1.POST("/areas", usercontroller.AuthUser(model.EDITOR), controller.AddArea)
		v1.GET("/products/:id/areas", usercontroller.AuthUser(model.EDITOR), controller.GetProductAreas)
		v1.PUT("/areas/:id", usercontroller.AuthUser(model.EDITOR), controller.UpdateArea)
		v1.DELETE("/areas/:id", usercontroller.AuthUser(model.EDITOR), controller.DeleteArea)

		v1.POST("/features", usercontroller.AuthUser(model.EDITOR), controller.AddFeature)
		v1.GET("/areas/:id/features", usercontroller.AuthUser(model.EDITOR), controller.GetAreaFeatures)
		v1.PUT("/features/:id", usercontroller.AuthUser(model.EDITOR), controller.UpdateFeature)
		v1.DELETE("/features/:id", usercontroller.AuthUser(model.EDITOR), controller.DeleteFeature)

		v1.POST("/tests", usercontroller.AuthUser(model.EDITOR), controller.AddTest)
		v1.GET("/tests", usercontroller.AuthUser(model.EDITOR), controller.GetAllTestForSuiteFile)
		v1.DELETE("/tests/:id", usercontroller.AuthUser(model.EDITOR), controller.DeleteTest)

		v1.POST("/expl-tests", usercontroller.AuthUser(model.CONSUMER), controller.AddExplTest)
		v1.GET("/expl-tests/area/:areaid", usercontroller.AuthUser(model.CONSUMER), controller.GetExplTestsForArea)
		v1.DELETE("/expl-tests/:id", usercontroller.AuthUser(model.EDITOR), controller.DeleteExplTest)

		v1.POST("/coverage/:id/upload-mocha-summary-report", AuthApi(), controller.UploadMochaSummaryReport)
		v1.GET("/coverage/:id/areas", usercontroller.AuthUser(model.CONSUMER), controller.GetAreaCoverage)
		v1.GET("/coverage/areas/:id/features", usercontroller.AuthUser(model.CONSUMER), controller.GetFeatureCoverage)
		v1.GET("/coverage/features/:id/tests", usercontroller.AuthUser(model.CONSUMER), controller.GetTestsCoverage)
		v1.GET("/coverage/products/:id/tests", usercontroller.AuthUser(model.CONSUMER), controller.GetProductTestsCoverage)

		// User related endpoints
		v1.POST("/auth/signin", usercontroller.Signin)
		v1.POST("/auth/refresh", usercontroller.RefreshToken)
	}

	port, found := os.LookupEnv("PORT")
	if found {
		router.Run("0.0.0.0:" + port)
	} else {
		router.Run("0.0.0.0:8080")
	}
}
