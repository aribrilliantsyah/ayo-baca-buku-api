// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "email": "fiber@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/users": {
            "get": {
                "description": "Get all users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get all users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ayo-baca-buku_app_models.User"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ayo-baca-buku_app_models.ReadingActivity": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "$ref": "#/definitions/gorm.DeletedAt"
                },
                "end_page": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "notes": {
                    "type": "string"
                },
                "pages_read": {
                    "type": "integer"
                },
                "reading_date": {
                    "type": "string"
                },
                "start_page": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "user_book": {
                    "$ref": "#/definitions/ayo-baca-buku_app_models.UserBook"
                },
                "user_book_id": {
                    "type": "integer"
                }
            }
        },
        "ayo-baca-buku_app_models.User": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "$ref": "#/definitions/gorm.DeletedAt"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "user_books": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ayo-baca-buku_app_models.UserBook"
                    }
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "ayo-baca-buku_app_models.UserBook": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "cover": {
                    "description": "URL atau path ke gambar cover",
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "current_page": {
                    "type": "integer"
                },
                "deleted_at": {
                    "$ref": "#/definitions/gorm.DeletedAt"
                },
                "end_date": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "motivation_read": {
                    "type": "string"
                },
                "publisher": {
                    "type": "string"
                },
                "reading_activities": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ayo-baca-buku_app_models.ReadingActivity"
                    }
                },
                "start_date": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "total_pages": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/ayo-baca-buku_app_models.User"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "gorm.DeletedAt": {
            "type": "object",
            "properties": {
                "time": {
                    "type": "string"
                },
                "valid": {
                    "description": "Valid is true if Time is not NULL",
                    "type": "boolean"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Ayo Baca Buku - API",
	Description:      "Ini adalah API - Ayo Baca Buku",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}