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
        "termsOfService": "https://www.bjungle.net/terms/",
        "contact": {
            "name": "API Support",
            "email": "luis.lucero@bjungle.net"
        },
        "license": {
            "name": "Software Owner",
            "url": "https://www.bjungle.net/terms/licenses"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/traceability/user-session/{userID}": {
            "get": {
                "description": "Método para obtener la trazabilidad registrada para el proceso de verificación de identidad",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Traceability"
                ],
                "summary": "Obtención de los datos de trazabilidad",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID del usuario",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/traceability.resTraceability"
                        }
                    }
                }
            }
        },
        "/api/v1/user/basic-information": {
            "post": {
                "description": "Método para el registro de los datos básicos de una persona",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Registro de información básica",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "request of validate user identity",
                        "name": "BasicInformation",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.requestValidateIdentity"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.resCreateUser"
                        }
                    }
                }
            }
        },
        "/api/v1/user/create": {
            "post": {
                "description": "Método para crear el usuario",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Creación de un usuario",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "request of validate user identity",
                        "name": "BasicInformation",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.requestValidateIdentity"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.responseAnny"
                        }
                    }
                }
            }
        },
        "/api/v1/user/data-pending": {
            "get": {
                "description": "Método para el obtener la cantidad de usuarios que no han cargado la información básica como la selfie, el documento de identidad y la información básica",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Obtiene la cantidad de usuarios que no cargaron información requerida",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.resGetUsersDataPending"
                        }
                    }
                }
            }
        },
        "/api/v1/user/upload-documents": {
            "post": {
                "description": "Método para cargar el documento de identidad",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Carga del documento de identidad",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Documento de identidad",
                        "name": "uploadDocument",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.reqUploadDocument"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.responseAnny"
                        }
                    }
                }
            }
        },
        "/api/v1/user/upload-selfie": {
            "post": {
                "description": "Método para cargar la selfie del usuario",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Carga de selfie del usuario",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Selfie del usuario",
                        "name": "UploadSelfie",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.reqUploadSelfie"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.responseAnny"
                        }
                    }
                }
            }
        },
        "/api/v1/user/user-session/{identifier}": {
            "get": {
                "description": "Método para el obtener la información del usuario en sesión por su email o id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Obtiene los datos registrados del usuario por su email o su id",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Identificador para la búsqueda del usuario",
                        "name": "identifier",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.resGetUserSession"
                        }
                    }
                }
            }
        },
        "/api/v1/user/users-lasted/{email}/{limit}/{offset}": {
            "get": {
                "description": "Método para el obtener los registros de los usuarios",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Obtiene los registros de usuarios",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Correo electrónico del usuario",
                        "name": "email",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Cantidad de registros por consulta",
                        "name": "limit",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Inicio del conteo de los registros por consulta",
                        "name": "offset",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.resGetUsersLasted"
                        }
                    }
                }
            }
        },
        "/api/v1/work/accept": {
            "post": {
                "description": "Método para aceptar la data registrada de un usuario por parte del administrador",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Work"
                ],
                "summary": "Acepta la información registrada",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Datos de solicitud para la aceptación",
                        "name": "ReqAccept",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/work.ReqAccept"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/work.resAnny"
                        }
                    }
                }
            }
        },
        "/api/v1/work/all": {
            "get": {
                "description": "Método para obtener la totalidad del trabajo registrado por lo usuarios",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Work"
                ],
                "summary": "Trae la totalidad del trabajo existente",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/work.resAllWork"
                        }
                    }
                }
            }
        },
        "/api/v1/work/refused": {
            "post": {
                "description": "Método para rechazar la data registrada de un usuario por parte del administrador",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Work"
                ],
                "summary": "Rechaza la información registrada",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Datos de solicitud para el rechazo",
                        "name": "ReqRefused",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/work.ReqRefused"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/work.resAnny"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "traceability.Traceability": {
            "type": "object",
            "properties": {
                "action": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "traceability.resTraceability": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/traceability.Traceability"
                    }
                },
                "error": {
                    "type": "boolean"
                },
                "msg": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "users.DataPending": {
            "type": "object",
            "properties": {
                "basic_information": {
                    "type": "integer"
                },
                "document": {
                    "type": "integer"
                },
                "selfie": {
                    "type": "integer"
                }
            }
        },
        "users.UserStatus": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "first_surname": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "second_name": {
                    "type": "string"
                },
                "second_surname": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "users.UserValidation": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer"
                },
                "back_document_img": {
                    "type": "string"
                },
                "birth_date": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "civil_status": {
                    "type": "string"
                },
                "country": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "department": {
                    "type": "string"
                },
                "document_number": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "expedition_date": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "first_surname": {
                    "type": "string"
                },
                "front_document_img": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "nationality": {
                    "type": "string"
                },
                "second_name": {
                    "type": "string"
                },
                "second_surname": {
                    "type": "string"
                },
                "selfie_img": {
                    "type": "string"
                },
                "transaction_id": {
                    "type": "string"
                },
                "type_document": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "users.Users": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer"
                },
                "birth_date": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "civil_status": {
                    "type": "string"
                },
                "country": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "department": {
                    "type": "string"
                },
                "document_number": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "expedition_date": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "first_surname": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "nationality": {
                    "type": "string"
                },
                "real_ip": {
                    "type": "string"
                },
                "second_name": {
                    "type": "string"
                },
                "second_surname": {
                    "type": "string"
                },
                "type_document": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "users.reqUploadDocument": {
            "type": "object",
            "properties": {
                "document_back_img": {
                    "type": "string"
                },
                "document_front_img": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "users.reqUploadSelfie": {
            "type": "object",
            "properties": {
                "selfie_img": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "users.requestValidateIdentity": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer"
                },
                "birth_date": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "civil_status": {
                    "type": "string"
                },
                "country": {
                    "type": "string"
                },
                "department": {
                    "type": "string"
                },
                "document_number": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "expedition_date": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "first_surname": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "nationality": {
                    "type": "string"
                },
                "second_name": {
                    "type": "string"
                },
                "second_surname": {
                    "type": "string"
                },
                "type_document": {
                    "type": "string"
                }
            }
        },
        "users.resCreateUser": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "$ref": "#/definitions/users.Users"
                },
                "error": {
                    "type": "boolean"
                },
                "msg": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "users.resGetUserSession": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "$ref": "#/definitions/users.UserValidation"
                },
                "error": {
                    "type": "boolean"
                },
                "msg": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "users.resGetUsersDataPending": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "$ref": "#/definitions/users.DataPending"
                },
                "error": {
                    "type": "boolean"
                },
                "msg": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "users.resGetUsersLasted": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/users.UserStatus"
                    }
                },
                "error": {
                    "type": "boolean"
                },
                "msg": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "users.responseAnny": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "string"
                },
                "error": {
                    "type": "boolean"
                },
                "msg": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "work.ReqAccept": {
            "type": "object",
            "properties": {
                "user_id": {
                    "type": "string"
                }
            }
        },
        "work.ReqRefused": {
            "type": "object",
            "properties": {
                "motivo": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "work.Status": {
            "type": "object",
            "properties": {
                "expired": {
                    "type": "integer"
                },
                "not_stated": {
                    "type": "integer"
                },
                "pending": {
                    "type": "integer"
                },
                "refused": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                },
                "valid": {
                    "type": "integer"
                }
            }
        },
        "work.resAllWork": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "$ref": "#/definitions/work.Status"
                },
                "error": {
                    "type": "boolean"
                },
                "msg": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "work.resAnny": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "string"
                },
                "error": {
                    "type": "boolean"
                },
                "msg": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        }
    },
    "tags": [
        {
            "description": "Métodos referentes al usuario",
            "name": "User"
        },
        {
            "description": "Métodos referentes a la trazabilidad",
            "name": "Traceability"
        },
        {
            "description": "Métodos referentes al trabajo registrado",
            "name": "Work"
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.4",
	Host:             "http://172.147.77.149:50050",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Check ID OnBoarding",
	Description:      "Api para OnBoarding y validación de identidad",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
