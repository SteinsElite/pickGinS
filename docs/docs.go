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
        "/api/v1/chart/profit": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "get the phased profit info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the duration to query-{7D,1M,1Y}",
                        "name": "range",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "the [(timestamp,profit)] in the time range"
                    }
                }
            }
        },
        "/api/v1/chart/ratio": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "get the ratio info",
                "responses": {
                    "200": {
                        "description": "amount of each asset in usd"
                    }
                }
            }
        },
        "/api/v1/chart/volume": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "get the total volume info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the duration to query-{7D,1M,1Y}",
                        "name": "range",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "the [(timestamp,volume)] in the time range"
                    }
                }
            }
        },
        "/api/v1/price_info/{coin}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "get the Coin Price info and trend",
                "parameters": [
                    {
                        "type": "string",
                        "description": "{BTC,ETH,USDT,HT,MDX}",
                        "name": "coin",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "the price trend of coin, {\"rate\": ..., \"trend\": ...}"
                    }
                }
            }
        },
        "/api/v1/transaction/{address}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "get the transaction info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user account address",
                        "name": "address",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "tag of the transaction-{deposit,withdraw,claimProfit}, if not specify, get all the category",
                        "name": "tag",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "index of page",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "size of each page",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "the transaction of the page"
                    },
                    "400": {
                        "description": "Invalid params"
                    },
                    "500": {
                        "description": "Server error"
                    }
                }
            }
        }
    },
    "definitions": {
        "notification.Notification": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "string"
                },
                "content": {
                    "type": "string"
                },
                "timeStamp": {
                    "type": "integer"
                },
                "title": {
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
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
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
	swag.Register(swag.Name, &s{})
}
