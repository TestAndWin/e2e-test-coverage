// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/areas": {
            "put": {
                "description": "Takes an area JSON and the Area ID and updates an area in the DB.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "area"
                ],
                "summary": "Updates an area",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Area ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Area JSON",
                        "name": "area",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Area"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Area"
                        }
                    }
                }
            },
            "post": {
                "description": "Takes an area JSON and stores it in DB. Return saved JSON.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "area"
                ],
                "summary": "Add a new area to a product",
                "parameters": [
                    {
                        "description": "Area JSON",
                        "name": "area",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Area"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Area"
                        }
                    }
                }
            }
        },
        "/api/v1/areas/{id}": {
            "delete": {
                "description": "Deletes the product area",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "area"
                ],
                "summary": "Deletes the product area",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Area ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/v1/areas/{id}/features": {
            "get": {
                "description": "Get all features for the specified area",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "feature"
                ],
                "summary": "Get all area features",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Area ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Feature"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/coverage/areas/{id}/features": {
            "get": {
                "description": "Get coverage for all area features. Only tests from the last 28 days are considered.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "coverage"
                ],
                "summary": "Get coverage for all area features.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Area ID",
                        "name": "product",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Feature"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/coverage/{id}/areas": {
            "get": {
                "description": "Get coverage for all product areas. Only tests from the last 28 days are considered.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "coverage"
                ],
                "summary": "Get coverage for all product areas.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Product ID",
                        "name": "product",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Area"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/expl-tests": {
            "post": {
                "description": "Takes a exploratory test JSON and stores it in DB. Return saved JSON.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "expl-test"
                ],
                "summary": "Add a new expl test",
                "parameters": [
                    {
                        "description": "Expl Test JSON",
                        "name": "expl-test",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ExplTest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.ExplTest"
                        }
                    }
                }
            }
        },
        "/api/v1/expl-tests/area/{areaid}": {
            "post": {
                "description": "Get all exploratory tests for the specified area for the last 28 days",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "expl-test"
                ],
                "summary": "Get all exploratory tests.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Area ID",
                        "name": "areaid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.ExplTest"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/expl-tests/{id}": {
            "delete": {
                "description": "Deletes the expl test",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "expl-test"
                ],
                "summary": "Deletes the expl test",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Test ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/v1/features": {
            "put": {
                "description": "Takes a feature JSON and feature ID and updates it in DB. Return saved JSON.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "feature"
                ],
                "summary": "Updates a feature",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Feature ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Feature JSON",
                        "name": "feature",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Feature"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Feature"
                        }
                    }
                }
            },
            "post": {
                "description": "Takes a feature JSON and stores it in DB. Return saved JSON.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "feature"
                ],
                "summary": "Add a new feature to an area",
                "parameters": [
                    {
                        "description": "Feature JSON",
                        "name": "feature",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Feature"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Feature"
                        }
                    }
                }
            }
        },
        "/api/v1/features/{id}": {
            "delete": {
                "description": "Deletes the product feature",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "feature"
                ],
                "summary": "Deletes the product feature",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Feature ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/v1/products": {
            "get": {
                "description": "Get all products",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "Get all products",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Product"
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "Takes a product JSON and product ID and updates it in DB. Return saved JSON.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "Updates a product",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Product JSON",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Product"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Product"
                        }
                    }
                }
            },
            "post": {
                "description": "Takes a product JSON and stores it in DB. Return saved JSON.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "Add a new product",
                "parameters": [
                    {
                        "description": "Product JSON",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Product"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Product"
                        }
                    }
                }
            }
        },
        "/api/v1/products/{id}": {
            "delete": {
                "description": "Deletes the product",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "Deletes the product",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/v1/products/{id}/areas": {
            "get": {
                "description": "Get all areas for the specified product",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "area"
                ],
                "summary": "Get all product areas",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Area"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/tests": {
            "get": {
                "description": "Get all tests for the specified suite and filename.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "test"
                ],
                "summary": "Get all tests for the specified suite and filename.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Suite name",
                        "name": "suite",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "File name",
                        "name": "file-name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Test"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Takes a test JSON and stores it in DB. Return saved JSON.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "test"
                ],
                "summary": "Add a new test",
                "parameters": [
                    {
                        "description": "Test JSON",
                        "name": "test",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Test"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Test"
                        }
                    }
                }
            }
        },
        "/api/v1/tests/{id}": {
            "delete": {
                "description": "Deletes the test",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "test"
                ],
                "summary": "Deletes the test",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Test ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/coverage/:id/upload-mocha-report": {
            "post": {
                "description": "Adds test result of a mocha report and returns the ID of the stored test.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "mocha"
                ],
                "summary": "Adds test result of a mocha report",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Mocha JSON",
                        "name": "test",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/coverage/features/:id/tests": {
            "get": {
                "description": "Get coverage for all tests of a feature for the last 28 days.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "coverage"
                ],
                "summary": "Get coverage for all tests of a feature.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Feature ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Test"
                            }
                        }
                    }
                }
            }
        },
        "/coverage/products/:id/tests": {
            "get": {
                "description": "Get coverage for all tests of a product for the last 28 days.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "coverage"
                ],
                "summary": "Get coverage for all tests of a product.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Test"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Area": {
            "type": "object",
            "properties": {
                "expl-rating": {
                    "type": "number"
                },
                "expl-tests": {
                    "type": "integer"
                },
                "failures": {
                    "type": "integer"
                },
                "first-total": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "passes": {
                    "type": "integer"
                },
                "pending": {
                    "type": "integer"
                },
                "product_id": {
                    "type": "integer"
                },
                "skipped": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "model.ExplTest": {
            "type": "object",
            "properties": {
                "area-id": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "rating": {
                    "type": "integer"
                },
                "summary": {
                    "type": "string"
                },
                "test-run": {
                    "type": "string"
                }
            }
        },
        "model.Feature": {
            "type": "object",
            "properties": {
                "area-id": {
                    "type": "integer"
                },
                "business-value": {
                    "type": "string"
                },
                "documentation": {
                    "type": "string"
                },
                "failures": {
                    "type": "integer"
                },
                "first-total": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "passes": {
                    "type": "integer"
                },
                "pending": {
                    "type": "integer"
                },
                "skipped": {
                    "type": "integer"
                },
                "tests": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Test"
                    }
                },
                "total": {
                    "type": "integer"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "model.Product": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "model.Test": {
            "type": "object",
            "properties": {
                "area-id": {
                    "type": "integer"
                },
                "failed-test-runs": {
                    "type": "integer"
                },
                "failures": {
                    "type": "integer"
                },
                "feature-id": {
                    "type": "integer"
                },
                "file-name": {
                    "type": "string"
                },
                "first-total": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "passes": {
                    "type": "integer"
                },
                "pending": {
                    "type": "integer"
                },
                "product-id": {
                    "type": "integer"
                },
                "skipped": {
                    "type": "integer"
                },
                "suite": {
                    "type": "string"
                },
                "test-run": {
                    "type": "string"
                },
                "total": {
                    "type": "integer"
                },
                "total-test-runs": {
                    "type": "integer"
                },
                "url": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "e2ecoverage",
	Description:      "API for e2e-coverage",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
