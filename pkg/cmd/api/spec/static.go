// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// Code generated by "makestatic"; DO NOT EDIT.

package static

var Files = map[string]string{
	"api.swagger.json": `{
  "swagger": "2.0",
  "info": {
    "title": "OpenPireix Project",
    "version": "0.0.1",
    "contact": {
      "name": "OpenPireix Project",
      "url": "https://openpitrix.io"
    }
  },
  "schemes": [
    "http"
  ],
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/v1/apps": {
      "get": {
        "summary": "describe apps with filter",
        "operationId": "DescribeApps",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/openpitrixDescribeAppsResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "app_id",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "name",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "repo_id",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "status",
            "in": "query",
            "required": false,
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          {
            "name": "limit",
            "in": "query",
            "required": false,
            "type": "integer",
            "format": "int64"
          },
          {
            "name": "offset",
            "in": "query",
            "required": false,
            "type": "integer",
            "format": "int64"
          }
        ],
        "tags": [
          "AppManager"
        ]
      },
      "delete": {
        "summary": "delete app",
        "operationId": "DeleteApp",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/openpitrixDeleteAppResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/openpitrixDeleteAppRequest"
            }
          }
        ],
        "tags": [
          "AppManager"
        ]
      },
      "post": {
        "summary": "create app",
        "operationId": "CreateApp",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/openpitrixCreateAppResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/openpitrixCreateAppRequest"
            }
          }
        ],
        "tags": [
          "AppManager"
        ]
      },
      "patch": {
        "summary": "modify app",
        "operationId": "ModifyApp",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/openpitrixModifyAppResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/openpitrixModifyAppRequest"
            }
          }
        ],
        "tags": [
          "AppManager"
        ]
      }
    },
    "/v1/credential_runtime_env": {
      "delete": {
        "summary": "detach runtime env",
        "operationId": "DetachCredentialFromRuntimeEnv",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/openpitrixDetachCredentialFromRuntimeEnvResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/openpitrixDetachCredentialFromRuntimeEnvRequset"
            }
          }
        ],
        "tags": [
          "RuntimeEnvManager"
        ]
      },
      "post": {
        "summary": "create runtime env",
        "operationId": "AttachCredentialToRuntimeEnv",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/openpitrixAttachCredentialToRuntimeEnvResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/openpitrixAttachCredentialToRuntimeEnvRequset"
            }
          }
        ],
        "tags": [
          "RuntimeEnvManager"
        ]
      }
    },
    "/v1/runtime_env_credentials": {
      "get": {
        "summary": "describe runtime env crendentials",
        "operationId": "DescribeRuntimeEnvCredentials",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/openpitrixDescribeRuntimeEnvCredentialsResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "runtime_env_credential_ids",
            "in": "query",
            "required": false,
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          {
            "name": "statuses",
            "in": "query",
            "required": false,
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          {
            "name": "search_word",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "owners",
            "in": "query",
            "required": false,
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          {
            "name": "limit.value",
            "description": "The uint32 value.",
            "in": "query",
            "required": false,
            "type": "integer",
            "format": "int64"
          },
          {
            "name": "offset.value",
            "description": "The uint32 value.",
            "in": "query",
            "required": false,
            "type": "integer",
            "format": "int64"
          },
          {
            "name": "verbose.value",
            "description": "The uint32 value.",
            "in": "query",
            "required": false,
            "type": "integer",
            "format": "int64"
          }
        ],
        "tags": [
          "RuntimeEnvManager"
        ]
      },
      "delete": {
        "summary": "modify runtime env credential",
        "operationId": "DeleteRuntimeEnvCredential",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/openpitrixDeleteRuntimeEnvCredentialResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/openpitrixDeleteRuntimeEnvCredentialRequset"
            }
          }
        ],
        "tags": [
          "RuntimeEnvManager"
        ]
      },
      "post": {
        "summary": "create runtime env credential",
        "operationId": "CreateRuntimeEnvCredential",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/openpitrixCreateRuntimeEnvCredentialResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/openpitrixCreateRuntimeEnvCredentialRequset"
            }
          }
        ],
        "tags": [
          "RuntimeEnvManager"
        ]
      },
      "patch": {
        "summary": "modify runtime env credential",
        "operationId": "ModifyRuntimeEnvCredential",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/openpitrixModifyRuntimeEnvCredentialResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/openpitrixModifyRuntimeEnvCredentialRequest"
            }
          }
        ],
        "tags": [
          "RuntimeEnvManager"
        ]
      }
    },
    "/v1/runtime_envs": {
      "get": {
        "summary": "describe runtime envs",
        "operationId": "DescribeRuntimeEnvs",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/openpitrixDescribeRuntimeEnvsResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "runtime_env_ids",
            "in": "query",
            "required": false,
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          {
            "name": "statuses",
            "in": "query",
            "required": false,
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          {
            "name": "search_word",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "selector",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "owners",
            "in": "query",
            "required": false,
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          {
            "name": "limit.value",
            "description": "The uint32 value.",
            "in": "query",
            "required": false,
            "type": "integer",
            "format": "int64"
          },
          {
            "name": "offset.value",
            "description": "The uint32 value.",
            "in": "query",
            "required": false,
            "type": "integer",
            "format": "int64"
          },
          {
            "name": "verbose.value",
            "description": "The uint32 value.",
            "in": "query",
            "required": false,
            "type": "integer",
            "format": "int64"
          }
        ],
        "tags": [
          "RuntimeEnvManager"
        ]
      },
      "delete": {
        "summary": "create runtime env",
        "operationId": "DeleteRuntimeEnv",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/openpitrixDeleteRuntimeEnvResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/openpitrixDeleteRuntimeEnvRequest"
            }
          }
        ],
        "tags": [
          "RuntimeEnvManager"
        ]
      },
      "post": {
        "summary": "create runtime env",
        "operationId": "CreateRuntimeEnv",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/openpitrixCreateRuntimeEnvResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/openpitrixCreateRuntimeEnvRequest"
            }
          }
        ],
        "tags": [
          "RuntimeEnvManager"
        ]
      },
      "patch": {
        "summary": "modify runtime env",
        "operationId": "ModifyRuntimeEnv",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/openpitrixModifyRuntimeEnvResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/openpitrixModifyRuntimeEnvRequest"
            }
          }
        ],
        "tags": [
          "RuntimeEnvManager"
        ]
      }
    }
  },
  "definitions": {
    "openpitrixApp": {
      "type": "object",
      "properties": {
        "app_id": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "repo_id": {
          "type": "string"
        },
        "description": {
          "type": "string"
        },
        "status": {
          "type": "string"
        },
        "home": {
          "type": "string"
        },
        "icon": {
          "type": "string"
        },
        "screenshots": {
          "type": "string"
        },
        "maintainers": {
          "type": "string"
        },
        "sources": {
          "type": "string"
        },
        "readme": {
          "type": "string"
        },
        "chart_name": {
          "type": "string"
        },
        "owner": {
          "type": "string"
        },
        "create_time": {
          "type": "string",
          "format": "date-time"
        },
        "status_time": {
          "type": "string",
          "format": "date-time"
        }
      }
    },
    "openpitrixCreateAppRequest": {
      "type": "object",
      "properties": {
        "_": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "repo_id": {
          "type": "string"
        },
        "owner": {
          "type": "string"
        },
        "chart_name": {
          "type": "string"
        },
        "description": {
          "type": "string"
        },
        "home": {
          "type": "string"
        },
        "icon": {
          "type": "string"
        },
        "screenshots": {
          "type": "string"
        },
        "maintainers": {
          "type": "string"
        },
        "sources": {
          "type": "string"
        },
        "readme": {
          "type": "string"
        }
      }
    },
    "openpitrixCreateAppResponse": {
      "type": "object",
      "properties": {
        "app": {
          "$ref": "#/definitions/openpitrixApp"
        }
      }
    },
    "openpitrixDeleteAppRequest": {
      "type": "object",
      "properties": {
        "app_id": {
          "type": "string"
        }
      }
    },
    "openpitrixDeleteAppResponse": {
      "type": "object",
      "properties": {
        "app": {
          "$ref": "#/definitions/openpitrixApp"
        }
      }
    },
    "openpitrixDescribeAppsResponse": {
      "type": "object",
      "properties": {
        "total_count": {
          "type": "integer",
          "format": "int64"
        },
        "app_set": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/openpitrixApp"
          }
        }
      }
    },
    "openpitrixModifyAppRequest": {
      "type": "object",
      "properties": {
        "app_id": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "repo_id": {
          "type": "string"
        },
        "owner": {
          "type": "string"
        },
        "chart_name": {
          "type": "string"
        },
        "description": {
          "type": "string"
        },
        "home": {
          "type": "string"
        },
        "icon": {
          "type": "string"
        },
        "screenshots": {
          "type": "string"
        },
        "maintainers": {
          "type": "string"
        },
        "sources": {
          "type": "string"
        },
        "readme": {
          "type": "string"
        }
      }
    },
    "openpitrixModifyAppResponse": {
      "type": "object",
      "properties": {
        "app": {
          "$ref": "#/definitions/openpitrixApp"
        }
      }
    },
    "protobufStringValue": {
      "type": "object",
      "properties": {
        "value": {
          "type": "string",
          "description": "The string value."
        }
      },
      "description": "Wrapper message for ` + "`" + `string` + "`" + `.\n\nThe JSON representation for ` + "`" + `StringValue` + "`" + ` is JSON string."
    },
    "openpitrixAttachCredentialToRuntimeEnvRequset": {
      "type": "object",
      "properties": {
        "runtime_env_credential_id": {
          "type": "string"
        },
        "runtime_env_id": {
          "type": "string"
        }
      }
    },
    "openpitrixAttachCredentialToRuntimeEnvResponse": {
      "type": "object",
      "properties": {
        "runtime_env_credential_id": {
          "type": "string"
        },
        "runtime_env_id": {
          "type": "string"
        }
      }
    },
    "openpitrixCreateRuntimeEnvCredentialRequset": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "description": {
          "type": "string"
        },
        "content": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        }
      }
    },
    "openpitrixCreateRuntimeEnvCredentialResponse": {
      "type": "object",
      "properties": {
        "runtime_env_credential": {
          "$ref": "#/definitions/openpitrixRuntimeEnvCredential"
        }
      }
    },
    "openpitrixCreateRuntimeEnvRequest": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "description": {
          "type": "string"
        },
        "labels": {
          "type": "string"
        },
        "runtime_env_url": {
          "type": "string"
        }
      }
    },
    "openpitrixCreateRuntimeEnvResponse": {
      "type": "object",
      "properties": {
        "runtime_env": {
          "$ref": "#/definitions/openpitrixRuntimeEnv"
        }
      }
    },
    "openpitrixDeleteRuntimeEnvCredentialRequset": {
      "type": "object",
      "properties": {
        "runtime_env_credential_id": {
          "type": "string"
        }
      }
    },
    "openpitrixDeleteRuntimeEnvCredentialResponse": {
      "type": "object",
      "properties": {
        "runtime_env_credential": {
          "$ref": "#/definitions/openpitrixRuntimeEnvCredential"
        }
      }
    },
    "openpitrixDeleteRuntimeEnvRequest": {
      "type": "object",
      "properties": {
        "runtime_env_id": {
          "type": "string"
        }
      }
    },
    "openpitrixDeleteRuntimeEnvResponse": {
      "type": "object",
      "properties": {
        "runtime_env": {
          "$ref": "#/definitions/openpitrixRuntimeEnv"
        }
      }
    },
    "openpitrixDescribeRuntimeEnvCredentialsResponse": {
      "type": "object",
      "properties": {
        "total_count": {
          "$ref": "#/definitions/protobufUInt32Value"
        },
        "runtime_env_credential_set": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/openpitrixRuntimeEnvCredential"
          }
        }
      }
    },
    "openpitrixDescribeRuntimeEnvsResponse": {
      "type": "object",
      "properties": {
        "total_count": {
          "$ref": "#/definitions/protobufUInt32Value"
        },
        "runtime_env_set": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/openpitrixRuntimeEnv"
          }
        }
      }
    },
    "openpitrixDetachCredentialFromRuntimeEnvRequset": {
      "type": "object",
      "properties": {
        "runtime_env_credential_id": {
          "type": "string"
        },
        "runtime_env_id": {
          "type": "string"
        }
      }
    },
    "openpitrixDetachCredentialFromRuntimeEnvResponse": {
      "type": "object",
      "properties": {
        "runtime_env_credential_id": {
          "type": "string"
        },
        "runtime_env_id": {
          "type": "string"
        }
      }
    },
    "openpitrixModifyRuntimeEnvCredentialRequest": {
      "type": "object",
      "properties": {
        "runtime_env_credential_id": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "description": {
          "type": "string"
        },
        "content": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        }
      }
    },
    "openpitrixModifyRuntimeEnvCredentialResponse": {
      "type": "object",
      "properties": {
        "runtime_env_credential": {
          "$ref": "#/definitions/openpitrixRuntimeEnvCredential"
        }
      }
    },
    "openpitrixModifyRuntimeEnvRequest": {
      "type": "object",
      "properties": {
        "runtime_env_id": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "description": {
          "type": "string"
        },
        "labels": {
          "type": "string"
        }
      }
    },
    "openpitrixModifyRuntimeEnvResponse": {
      "type": "object",
      "properties": {
        "runtime_env": {
          "$ref": "#/definitions/openpitrixRuntimeEnv"
        }
      }
    },
    "openpitrixRuntimeEnv": {
      "type": "object",
      "properties": {
        "runtime_env_id": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "description": {
          "type": "string"
        },
        "runtime_env_url": {
          "type": "string"
        },
        "runtime_env_credential_id": {
          "type": "string"
        },
        "labels": {
          "type": "string"
        },
        "owner": {
          "type": "string"
        },
        "status": {
          "type": "string"
        },
        "create_time": {
          "type": "string",
          "format": "date-time"
        },
        "status_time": {
          "type": "string",
          "format": "date-time"
        }
      }
    },
    "openpitrixRuntimeEnvCredential": {
      "type": "object",
      "properties": {
        "runtime_env_credential_id": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "description": {
          "type": "string"
        },
        "content": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "owner": {
          "type": "string"
        },
        "runtime_env_ids": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "status": {
          "type": "string"
        },
        "create_time": {
          "type": "string",
          "format": "date-time"
        },
        "status_time": {
          "type": "string",
          "format": "date-time"
        }
      }
    },
    "protobufUInt32Value": {
      "type": "object",
      "properties": {
        "value": {
          "type": "integer",
          "format": "int64",
          "description": "The uint32 value."
        }
      },
      "description": "Wrapper message for ` + "`" + `uint32` + "`" + `.\n\nThe JSON representation for ` + "`" + `UInt32Value` + "`" + ` is JSON number."
    }
  }
}
`,

	"index.html": `<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Swagger UI</title>
  <link href="https://fonts.googleapis.com/css?family=Open+Sans:400,700|Source+Code+Pro:300,600|Titillium+Web:400,600,700" rel="stylesheet">
  <link rel="stylesheet" type="text/css" href="./swagger-ui.css" >
  <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
  <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
  <style>
    html
    {
      box-sizing: border-box;
      overflow: -moz-scrollbars-vertical;
      overflow-y: scroll;
    }
    *,
    *:before,
    *:after
    {
      box-sizing: inherit;
    }

    body {
      margin:0;
      background: #fafafa;
    }
  </style>
</head>

<body>

<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" style="position:absolute;width:0;height:0">
  <defs>
    <symbol viewBox="0 0 20 20" id="unlocked">
          <path d="M15.8 8H14V5.6C14 2.703 12.665 1 10 1 7.334 1 6 2.703 6 5.6V6h2v-.801C8 3.754 8.797 3 10 3c1.203 0 2 .754 2 2.199V8H4c-.553 0-1 .646-1 1.199V17c0 .549.428 1.139.951 1.307l1.197.387C5.672 18.861 6.55 19 7.1 19h5.8c.549 0 1.428-.139 1.951-.307l1.196-.387c.524-.167.953-.757.953-1.306V9.199C17 8.646 16.352 8 15.8 8z"></path>
    </symbol>

    <symbol viewBox="0 0 20 20" id="locked">
      <path d="M15.8 8H14V5.6C14 2.703 12.665 1 10 1 7.334 1 6 2.703 6 5.6V8H4c-.553 0-1 .646-1 1.199V17c0 .549.428 1.139.951 1.307l1.197.387C5.672 18.861 6.55 19 7.1 19h5.8c.549 0 1.428-.139 1.951-.307l1.196-.387c.524-.167.953-.757.953-1.306V9.199C17 8.646 16.352 8 15.8 8zM12 8H8V5.199C8 3.754 8.797 3 10 3c1.203 0 2 .754 2 2.199V8z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="close">
      <path d="M14.348 14.849c-.469.469-1.229.469-1.697 0L10 11.819l-2.651 3.029c-.469.469-1.229.469-1.697 0-.469-.469-.469-1.229 0-1.697l2.758-3.15-2.759-3.152c-.469-.469-.469-1.228 0-1.697.469-.469 1.228-.469 1.697 0L10 8.183l2.651-3.031c.469-.469 1.228-.469 1.697 0 .469.469.469 1.229 0 1.697l-2.758 3.152 2.758 3.15c.469.469.469 1.229 0 1.698z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="large-arrow">
      <path d="M13.25 10L6.109 2.58c-.268-.27-.268-.707 0-.979.268-.27.701-.27.969 0l7.83 7.908c.268.271.268.709 0 .979l-7.83 7.908c-.268.271-.701.27-.969 0-.268-.269-.268-.707 0-.979L13.25 10z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="large-arrow-down">
      <path d="M17.418 6.109c.272-.268.709-.268.979 0s.271.701 0 .969l-7.908 7.83c-.27.268-.707.268-.979 0l-7.908-7.83c-.27-.268-.27-.701 0-.969.271-.268.709-.268.979 0L10 13.25l7.418-7.141z"/>
    </symbol>


    <symbol viewBox="0 0 24 24" id="jump-to">
      <path d="M19 7v4H5.83l3.58-3.59L8 6l-6 6 6 6 1.41-1.41L5.83 13H21V7z"/>
    </symbol>

    <symbol viewBox="0 0 24 24" id="expand">
      <path d="M10 18h4v-2h-4v2zM3 6v2h18V6H3zm3 7h12v-2H6v2z"/>
    </symbol>

  </defs>
</svg>

<div id="swagger-ui"></div>

<script src="./swagger-ui-bundle.js"> </script>
<script src="./swagger-ui-standalone-preset.js"> </script>
<script>
window.onload = function() {

  // Build a system
  const ui = SwaggerUIBundle({
    urls: [
      { name:"Api", url:"/swagger-ui/api.swagger.json" }
    ],
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  })

  window.ui = ui
}
</script>
</body>

</html>
`,
}
