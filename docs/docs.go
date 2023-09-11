// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/api/v1/client": {
            "post": {
                "description": "Método para crear el cliente en el sistema",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "summary": "Crea el cliente en el sistema",
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
                        "description": "Datos del cliente a crear",
                        "name": "Client",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/clients.Client"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/clients.ResAnny"
                        }
                    }
                }
            }
        },
        "/api/v1/client/{nit}": {
            "get": {
                "description": "Método para obtener la información del cliente de CheckID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "summary": "Obtiene la data del cliente",
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
                        "description": "NIT del cliente",
                        "name": "nit",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/clients.ResClient"
                        }
                    }
                }
            }
        },
        "/api/v1/onboarding/": {
            "post": {
                "description": "Método que permite iniciar el enrolamiento de un usuario que puede ser desde un tercero o desde el mismo sistema",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Onboarding"
                ],
                "summary": "Método que permite iniciar el enrolamiento de un usuario",
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
                        "description": "Datos para el enrolamiento del usuario",
                        "name": "resCreateOnboarding",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/onboarding.resCreateOnboarding"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/onboarding.resCreateOnboarding"
                        }
                    }
                }
            }
        },
        "/api/v1/onboarding/process": {
            "post": {
                "description": "Método que permite terminar el enrolamiento de un usuario que ha sido validado desde OnlyOne",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Onboarding"
                ],
                "summary": "Método que permite terminar el enrolamiento de un usuario",
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
                        "description": "Datos para validar el enrolamiento del usuario",
                        "name": "RequestProcessOnboarding",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/onboarding.RequestProcessOnboarding"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/onboarding.ResProcessOnboarding"
                        }
                    }
                }
            }
        },
        "/api/v1/onboarding/validate_identity": {
            "post": {
                "description": "Método que permite finalizar la validación de identidad de un usuario por la aplicación de OnlyOne",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Onboarding"
                ],
                "summary": "Método que permite finalizar la validación de identidad de un usuario",
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
                        "description": "Datos para validar la identidad del usuario",
                        "name": "RequestValidationIdentity",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/onboarding.RequestValidationIdentity"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/onboarding.ResProcessOnboarding"
                        }
                    }
                }
            }
        },
        "/api/v1/traceability": {
            "get": {
                "description": "Método para obtención de los datos de trazabilidad de un usuario por su id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Traceability"
                ],
                "summary": "Obtención de los datos de trazabilidad de un usuario por su id",
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
        "/api/v1/user/create": {
            "post": {
                "description": "Metodo que permite la creación de un usuario",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Metodo que permite la creación de un usuario",
                "parameters": [
                    {
                        "description": "request of validate user identity",
                        "name": "BasicInformation",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.RequestCreateUser"
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
        "/api/v1/user/file": {
            "get": {
                "description": "Método que permite validar si ha terminado la validación de identidad de un usuario",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Permite validar si ha terminado la validación de identidad de un usuario",
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
                        "description": "Id del archivo",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.ResponseGetUserFile"
                        }
                    }
                }
            }
        },
        "/api/v1/user/finish-onboarding": {
            "get": {
                "description": "Método que permite validar si se ha finalizado el proceso de enrolamiento de un usuario",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Permite validar si ha terminado el enrolamiento de un usuario",
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
                            "$ref": "#/definitions/users.responseFinishOnboarding"
                        }
                    }
                }
            }
        },
        "/api/v1/user/finish-validation": {
            "get": {
                "description": "Método que permite validar si ha terminado la validación de identidad de un usuario",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Permite validar si ha terminado la validación de identidad de un usuario",
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
                            "$ref": "#/definitions/users.responseFinishOnboarding"
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
        "/api/v1/user/user-session": {
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
        "/api/v1/user/validate": {
            "get": {
                "description": "Método para verificar si el usuario ha validado su identidad",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Verifica si el usuario ha validado su identidad",
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
                            "$ref": "#/definitions/users.responseAnny"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "check-id-api_api_handlers_onboarding.Onboarding": {
            "type": "object",
            "properties": {
                "method": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "clients.Client": {
            "type": "object",
            "properties": {
                "banner": {
                    "type": "string"
                },
                "full_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "logo_small": {
                    "type": "string"
                },
                "main_color": {
                    "type": "string"
                },
                "nit": {
                    "type": "string"
                },
                "second_color": {
                    "type": "string"
                },
                "url_api": {
                    "type": "string"
                },
                "url_redirect": {
                    "type": "string"
                }
            }
        },
        "clients.ResAnny": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
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
        "clients.ResClient": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "$ref": "#/definitions/clients.Client"
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
        "onboarding.RequestProcessOnboarding": {
            "type": "object",
            "properties": {
                "document_back": {
                    "type": "string"
                },
                "document_front": {
                    "type": "string"
                },
                "onboarding": {
                    "type": "string"
                },
                "selfie": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "onboarding.RequestValidationIdentity": {
            "type": "object",
            "properties": {
                "selfie": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                },
                "validation_id": {
                    "type": "integer"
                }
            }
        },
        "onboarding.ResProcessOnboarding": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
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
        "onboarding.resCreateOnboarding": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "$ref": "#/definitions/check-id-api_api_handlers_onboarding.Onboarding"
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
        "users.RequestCreateUser": {
            "type": "object",
            "properties": {
                "cellphone": {
                    "type": "string"
                },
                "document_number": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "users.ResponseGetUserFile": {
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
        "users.User": {
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
                    "type": "string"
                },
                "email": {
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
                "role": {
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
        "users.resGetUserSession": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "$ref": "#/definitions/users.User"
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
                "data": {},
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
        "users.responseFinishOnboarding": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "boolean"
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
            "description": "Métodos referentes al cliente",
            "name": "Client"
        },
        {
            "description": "Métodos referentes al enrolamiento del usuario",
            "name": "Onboarding"
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.4",
	Host:             "172.147.77.149:50050",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Check ID OnBoarding",
	Description:      "Api para OnBoarding y validación de identidad",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
