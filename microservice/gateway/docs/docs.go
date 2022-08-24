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
        "/auth/login/student": {
            "post": {
                "description": "login the student-forum",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "学生登录 api",
                "parameters": [
                    {
                        "description": "login_request",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/StudentLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/StudentLoginResponse"
                        }
                    }
                }
            }
        },
        "/auth/login/team": {
            "post": {
                "description": "login the team-forum",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "团队登录 api",
                "parameters": [
                    {
                        "description": "login_request",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/TeamLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/TeamLoginResponse"
                        }
                    }
                }
            }
        },
        "/chat": {
            "get": {
                "description": "该用户发送信息前先获取自己的uuid，并放入query(id=?)，有效期24h",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "获取该用户的uuid",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/chat.Id"
                        }
                    }
                }
            }
        },
        "/chat/ws": {
            "get": {
                "description": "建立 WebSocket 连接",
                "tags": [
                    "chat"
                ],
                "summary": "WebSocket",
                "parameters": [
                    {
                        "type": "string",
                        "description": "uuid",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/comment": {
            "post": {
                "description": "typeName: first-level -\u003e 一级评论; second-level -\u003e 其它级",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "comment"
                ],
                "summary": "创建评论 api",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "create_comment_request",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/comment.CreateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/comment/{comment_id}": {
            "get": {
                "description": "typeName: first-level -\u003e 一级评论; second-level -\u003e 其它级",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "comment"
                ],
                "summary": "获取评论 api",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "comment_id",
                        "name": "comment_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/comment.Comment"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "comment"
                ],
                "summary": "删除评论 api",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "comment_id",
                        "name": "comment_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/like": {
            "post": {
                "description": "TypeName: post or comment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "like"
                ],
                "summary": "点赞/取消点赞 api",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "like_request",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/like.Item"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/like/list": {
            "get": {
                "description": "TypeName: post or comment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "like"
                ],
                "summary": "获取用户点赞列表 api",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/like.ListResponse"
                        }
                    }
                }
            }
        },
        "/post": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "post"
                ],
                "summary": "修改帖子信息 api",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "update_info_request",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/post.UpdateInfoRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            },
            "post": {
                "description": "type_name : normal -\u003e 团队外; muxi -\u003e 团队内 (type_name暂时均填normal); 主帖的main_post_id不填或为0",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "post"
                ],
                "summary": "创建帖子 api",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "create_post_request",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/post.CreateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/post/list/{type_name}/sub/{main_post_id}": {
            "get": {
                "description": "type_name : normal -\u003e 团队外; muxi -\u003e 团队内 (type_name暂时均填normal); 根据 main_post_id 获取主帖的从帖list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "post"
                ],
                "summary": "list 从帖 api",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "last_id",
                        "name": "last_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "type_name",
                        "name": "type_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "main_post_id",
                        "name": "main_post_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/post.Post"
                            }
                        }
                    }
                }
            }
        },
        "/post/list/{type_name}/{category_id}": {
            "get": {
                "description": "type_name : normal -\u003e 团队外; muxi -\u003e 团队内 (type_name暂时均填normal); 根据category获取主帖list(前端实现category的映射)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "post"
                ],
                "summary": "list 主帖 api",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "last_id",
                        "name": "last_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "type_name",
                        "name": "type_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "category_id",
                        "name": "category_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/post.Post"
                            }
                        }
                    }
                }
            }
        },
        "/post/{post_id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "post"
                ],
                "summary": "获取帖子 api",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "post_id",
                        "name": "post_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/post.Post"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "post"
                ],
                "summary": "删除帖子 api",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "post_id",
                        "name": "post_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/user": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "修改用户个人信息 api",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "update_info_request",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UpdateInfoRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    }
                }
            }
        },
        "/user/list/{group_id}/{team_id}": {
            "get": {
                "description": "通过 group 和 team 获取 user_list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "list user api",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "last_id",
                        "name": "last_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "group_id",
                        "name": "group_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "team_id",
                        "name": "team_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ListResponse"
                        }
                    }
                }
            }
        },
        "/user/myprofile": {
            "get": {
                "description": "获取 my 完整 user 信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "get 我的 profile api",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/UserProfile"
                        }
                    }
                }
            }
        },
        "/user/profile/{id}": {
            "get": {
                "description": "通过 userId 获取完整 user 信息（权限: Normal-普通学生用户; NormalAdmin-学生管理员; Muxi-团队成员; MuxiAdmin-团队管理员; SuperAdmin-超级管理员）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "获取用户 profile api",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/UserProfile"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ListResponse": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/user"
                    }
                }
            }
        },
        "Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "StudentLoginRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "student_id": {
                    "type": "string"
                }
            }
        },
        "StudentLoginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "TeamLoginRequest": {
            "type": "object",
            "properties": {
                "oauth_code": {
                    "type": "string"
                }
            }
        },
        "TeamLoginResponse": {
            "type": "object",
            "properties": {
                "redirect_url": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "UpdateInfoRequest": {
            "type": "object",
            "properties": {
                "avatar_url": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "UserProfile": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
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
                "role": {
                    "type": "string"
                }
            }
        },
        "chat.Id": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "comment.Comment": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "create_time": {
                    "type": "string"
                },
                "creator_avatar": {
                    "type": "string"
                },
                "creator_id": {
                    "type": "integer"
                },
                "creator_name": {
                    "type": "string"
                },
                "father_id": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "is_liked": {
                    "type": "boolean"
                },
                "like_num": {
                    "type": "integer"
                },
                "type_name": {
                    "type": "string"
                }
            }
        },
        "comment.CreateRequest": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "father_id": {
                    "type": "integer"
                },
                "post_id": {
                    "type": "integer"
                },
                "type_name": {
                    "type": "string"
                }
            }
        },
        "like.Item": {
            "type": "object",
            "properties": {
                "target_id": {
                    "type": "integer"
                },
                "type_name": {
                    "type": "string"
                }
            }
        },
        "like.ListResponse": {
            "type": "object"
        },
        "post.CreateRequest": {
            "type": "object",
            "properties": {
                "category_id": {
                    "type": "integer"
                },
                "content": {
                    "type": "string"
                },
                "main_post_id": {
                    "type": "integer"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                },
                "type_name": {
                    "type": "string"
                }
            }
        },
        "post.Post": {
            "type": "object",
            "properties": {
                "category_id": {
                    "type": "integer"
                },
                "comment_num": {
                    "type": "integer"
                },
                "comments": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/comment.Comment"
                    }
                },
                "content": {
                    "type": "string"
                },
                "creator_avatar": {
                    "type": "string"
                },
                "creator_id": {
                    "type": "integer"
                },
                "creator_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "is_favorite": {
                    "type": "boolean"
                },
                "is_liked": {
                    "type": "boolean"
                },
                "like_num": {
                    "type": "integer"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "time": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "post.UpdateInfoRequest": {
            "type": "object",
            "properties": {
                "category_id": {
                    "type": "integer"
                },
                "content": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "user": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
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
                "role": {
                    "type": "string"
                }
            }
        }
    },
    "tags": [
        {
            "description": "用户服务",
            "name": "user"
        },
        {
            "description": "认证服务",
            "name": "auth"
        },
        {
            "description": "聊天服务",
            "name": "chat"
        },
        {
            "description": "帖子服务",
            "name": "post"
        }
    ]
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
	Host:        "forum.muxixyz.com",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "forum-gateway",
	Description: "The gateway of forum",
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
