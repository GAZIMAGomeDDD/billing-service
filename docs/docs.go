// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
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
        "/changeBalance": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "--",
                        "name": "changeBalance",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ChangeBalanceQuery"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "not enough money",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/getBalance": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "--",
                        "name": "getBalance",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.GetBalanceQuery"
                        }
                    },
                    {
                        "type": "string",
                        "description": "currency",
                        "name": "currency",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.GetBalanceResponse"
                        }
                    },
                    "404": {
                        "description": "user not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/getTransaction": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "--",
                        "name": "getTransaction",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.GetTransaction"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Transaction"
                        }
                    }
                }
            }
        },
        "/listOfTransactions": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "--",
                        "name": "listOfTransactions",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ListOfTransactionsQuery"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "page",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "sort",
                        "name": "sort",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Transaction"
                            }
                        }
                    }
                }
            }
        },
        "/moneyTransfer": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "--",
                        "name": "moneyTransfer",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.MoneyTransfer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "not enough money",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "user not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.ChangeBalanceQuery": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "b91a95a4-078f-4afd-b11c-4850eb65e784"
                },
                "money": {
                    "type": "number",
                    "example": 99.99
                }
            }
        },
        "model.ChangeBalanceResponse": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number",
                    "example": 99.99
                },
                "id": {
                    "type": "string",
                    "example": "b91a95a4-078f-4afd-b11c-4850eb65e784"
                }
            }
        },
        "model.GetBalanceQuery": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "b91a95a4-078f-4afd-b11c-4850eb65e784"
                }
            }
        },
        "model.GetBalanceResponse": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number",
                    "example": 99.99
                },
                "id": {
                    "type": "string",
                    "example": "b91a95a4-078f-4afd-b11c-4850eb65e784"
                }
            }
        },
        "model.GetTransaction": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "b91a95a4-078f-4afd-b11c-4850eb65e784"
                }
            }
        },
        "model.ListOfTransactionsQuery": {
            "type": "object",
            "properties": {
                "limit": {
                    "type": "integer",
                    "example": 5
                },
                "user_id": {
                    "type": "string",
                    "example": "b91a95a4-078f-4afd-b11c-4850eb65e784"
                }
            }
        },
        "model.MoneyTransfer": {
            "type": "object",
            "properties": {
                "from_id": {
                    "type": "string",
                    "example": "b81a95a4-078f-5dfd-b11c-4850eb35e785"
                },
                "money": {
                    "type": "number",
                    "example": 99.99
                },
                "to_id": {
                    "type": "string",
                    "example": "b91a95a4-078f-4afd-b11c-4850eb65e784"
                }
            }
        },
        "model.Transaction": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "from_id": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "method": {
                    "type": "string"
                },
                "money": {
                    "type": "number"
                },
                "to_id": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Billing Service API",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
