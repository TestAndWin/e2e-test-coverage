/*
Copyright (c) 2022-2023, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/TestAndWin/e2e-coverage/coverage/controller"
	"github.com/TestAndWin/e2e-coverage/ui"
	usercontroller "github.com/TestAndWin/e2e-coverage/user/controller"
	"github.com/TestAndWin/e2e-coverage/user/model"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/crypto/acme/autocert"
)

func HandleIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		html := ui.MustAsset("index.html")
		c.Data(200, "text/html; charset=UTF-8", html)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.GetHeader("origin"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func HandleRequest() {
	//gin.SetMode(gin.ReleaseMode)
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
		// Product ... config
		v1.POST("/products", usercontroller.AuthUser(model.MAINTAINER), controller.AddProduct)
		v1.GET("/products", usercontroller.AuthUser(model.MAINTAINER), controller.GetProducts)
		v1.PUT("/products/:id", usercontroller.AuthUser(model.MAINTAINER), controller.UpdateProduct)
		v1.DELETE("/products/:id", usercontroller.AuthUser(model.MAINTAINER), controller.DeleteProduct)

		v1.POST("/areas", usercontroller.AuthUser(model.MAINTAINER), controller.AddArea)
		v1.GET("/products/:id/areas", usercontroller.AuthUser(model.MAINTAINER), controller.GetProductAreas)
		v1.PUT("/areas/:id", usercontroller.AuthUser(model.MAINTAINER), controller.UpdateArea)
		v1.DELETE("/areas/:id", usercontroller.AuthUser(model.MAINTAINER), controller.DeleteArea)

		v1.POST("/features", usercontroller.AuthUser(model.MAINTAINER), controller.AddFeature)
		v1.GET("/areas/:id/features", usercontroller.AuthUser(model.MAINTAINER), controller.GetAreaFeatures)
		v1.PUT("/features/:id", usercontroller.AuthUser(model.MAINTAINER), controller.UpdateFeature)
		v1.DELETE("/features/:id", usercontroller.AuthUser(model.MAINTAINER), controller.DeleteFeature)

		v1.GET("/tests", usercontroller.AuthUser(model.MAINTAINER), controller.GetAllTestForSuiteFile)
		v1.DELETE("/tests/:id", usercontroller.AuthUser(model.MAINTAINER), controller.DeleteTest)

		// Expl Testing
		v1.POST("/expl-tests", usercontroller.AuthUser(model.TESTER), controller.AddExplTest)
		v1.GET("/expl-tests/area/:areaid", usercontroller.AuthUser(model.TESTER), controller.GetExplTestsForArea)
		v1.DELETE("/expl-tests/:id", usercontroller.AuthUser(model.MAINTAINER), controller.DeleteExplTest)

		// Test Coverage
		v1.POST("/coverage/:id/upload-mocha-summary-report", usercontroller.AuthApi(), controller.UploadMochaSummaryReport)
		v1.GET("/coverage/:id/areas", usercontroller.AuthUser(model.TESTER), controller.GetAreaCoverage)
		v1.GET("/coverage/areas/:id/features", usercontroller.AuthUser(model.TESTER), controller.GetFeatureCoverage)
		v1.GET("/coverage/features/:id/tests", usercontroller.AuthUser(model.TESTER), controller.GetTestsCoverage)
		v1.GET("/coverage/products/:id/tests", usercontroller.AuthUser(model.TESTER), controller.GetProductTestsCoverage)

		// User + admin related endpoints
		v1.POST("/auth/login", usercontroller.Login)
		v1.GET("/users", usercontroller.AuthUser(model.ADMIN), usercontroller.GetUser)
		v1.POST("/users", usercontroller.AuthUser(model.ADMIN), usercontroller.CreateUser)
		v1.PUT("/users/:id", usercontroller.AuthUser(model.ADMIN), usercontroller.UpdateUser)
		v1.DELETE("/users/:id", usercontroller.AuthUser(model.ADMIN), usercontroller.DeleteUser)
		v1.PUT("/users/change-pwd", usercontroller.AuthUser(""), usercontroller.ChangePassword)
		v1.POST("users/generate-api-key", usercontroller.AuthUser(model.ADMIN), usercontroller.GenerateApiKey)
	}

	_, devMode := os.LookupEnv("DEV")
	hostName, hostSet := os.LookupEnv("HOST")
	if devMode || !hostSet {
		fmt.Println("Start in DEV mode")
		router.Run("0.0.0.0:8080")
	} else {
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(hostName),
			Cache:      autocert.DirCache("/var/www/.cache"),
		}
		fmt.Println("Start in Prod mode, host white list: ", hostName)
		err := autotls.RunWithManager(router, &m)
		if err != nil {
			fmt.Println("Could not start in Prod mode: ", err)
		}
	}
}
