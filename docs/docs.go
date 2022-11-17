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
        "termsOfService": "There are no terms of service. We accept no responsibility for your ignorance.",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/lookup/{srv}/{usr}": {
            "get": {
                "description": "Gets the users account information by their Username and Server",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "accounts"
                ],
                "summary": "Show an account",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Riot Server",
                        "name": "srv",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "usr",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.LookupResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/match/{srv}/{usr}": {
            "get": {
                "description": "Show the past 10 matches",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "accounts"
                ],
                "summary": "Show recent matches",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Riot Server",
                        "name": "srv",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "usr",
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
                                "$ref": "#/definitions/models.MatchDataResp"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Ping the API service",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "utility"
                ],
                "summary": "Pings the API service to ensure that it is active",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "404": {
                        "description": "Not Found"
                    }
                }
            }
        },
        "/user": {
            "post": {
                "description": "Creates and stores the users data to be used when executing commands/api calls.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "accounts"
                ],
                "summary": "Create an account",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UserDB"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ChampionMasteryResp": {
            "type": "object",
            "properties": {
                "championId": {
                    "type": "integer"
                },
                "championLevel": {
                    "type": "integer"
                },
                "championName": {
                    "type": "string"
                },
                "championPoints": {
                    "type": "integer"
                }
            }
        },
        "models.LookupResponse": {
            "type": "object",
            "properties": {
                "champions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.ChampionMasteryResp"
                    }
                },
                "level": {
                    "type": "integer"
                },
                "losses": {
                    "type": "integer"
                },
                "profileIconId": {
                    "type": "integer"
                },
                "rank": {
                    "type": "string"
                },
                "tier": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                },
                "wins": {
                    "type": "integer"
                }
            }
        },
        "models.MatchDataResp": {
            "type": "object",
            "properties": {
                "info": {
                    "type": "object",
                    "properties": {
                        "gameMode": {
                            "type": "string"
                        },
                        "participants": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "properties": {
                                    "assists": {
                                        "type": "integer"
                                    },
                                    "championName": {
                                        "type": "string"
                                    },
                                    "deaths": {
                                        "type": "integer"
                                    },
                                    "kills": {
                                        "type": "integer"
                                    }
                                }
                            }
                        },
                        "teams": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "properties": {
                                    "win": {
                                        "type": "boolean"
                                    }
                                }
                            }
                        }
                    }
                }
            }
        },
        "models.UserDB": {
            "type": "object",
            "properties": {
                "RankedTier": {
                    "type": "string"
                },
                "discordid": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "puuid": {
                    "type": "string"
                },
                "server": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "api.LoLQueue.com",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "LoLQueue API",
	Description:      "This is the documentation for the LoLQueue Api Service",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
