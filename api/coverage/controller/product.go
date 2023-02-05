/*
Copyright (c) 2022-2023, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/TestAndWin/e2e-coverage/coverage/model"
	"github.com/gin-gonic/gin"
)

// AddProduct godoc
// @Summary      Add a new product
// @Description  Takes a product JSON and stores it in DB. Return saved JSON.
// @Tags         product
// @Produce      json
// @Param        product  body      model.Product  true  "Product JSON"
// @Success      201  {object}  model.Product
// @Router       /api/v1/products [POST]
func AddProduct(c *gin.Context) {
	var p model.Product
	if err := c.BindJSON(&p); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		id, err := repo.InsertProduct(p)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		} else {
			p.Id = id
			c.JSON(http.StatusCreated, p)
		}
	}
}

// GetProducts godoc
// @Summary      Get all products
// @Description  Get all products
// @Tags         product
// @Produce      json
// @Success      200  {array}  model.Product
// @Router       /api/v1/products [get]
func GetProducts(c *gin.Context) {
	p, err := repo.GetAllProducts()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		c.JSON(http.StatusOK, p)
	}
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Takes a product JSON and product ID and updates it in DB. Return saved JSON.
// @Tags         product
// @Param        id       path      int            true  "Product ID"
// @Param        product  body      model.Product  true  "Product JSON"
// @Produce      json
// @Success      200  {object}  model.Product
// @Router       /api/v1/products/{id} [PUT]
func UpdateProduct(c *gin.Context) {
	var p model.Product
	if err := c.BindJSON(&p); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		p.Id, _ = strconv.ParseInt(c.Param("id"), 0, 64)
		_, err := repo.UpdateProduct(p)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		} else {
			c.JSON(http.StatusOK, p)
		}
	}
}

// DeleteProduct godoc
// @Summary      Delete the product
// @Description  Delete the product
// @Tags         product
// @Produce      json
// @Param        id    path      int     true  "Product ID"
// @Success      200
// @Router       /api/v1/products/{id} [DELETE]
func DeleteProduct(c *gin.Context) {
	_, err := repo.DeleteProduct(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
