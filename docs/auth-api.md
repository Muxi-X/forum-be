# Auth API 接口文档

Base URL: `http://{host}/api/v1/auth`

公共响应体格式:

```json
{
  "code": 0,
  "message": "OK",
  "data": {}
}
```

`code` 为 `0` 表示成功，非 `0` 表示失败。`message` 为错误描述，`data` 为具体返回数据。

---

## 1. 注册

```
POST /api/v1/auth/register
Content-Type: application/json
```

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名，同时用作邮箱和昵称 |
| password | string | 是 | 密码 |

### 请求示例

```json
{
  "username": "user@example.com",
  "password": "123456"
}
```

### 成功响应

```json
{
  "code": 0,
  "message": "OK",
  "data": null
}
```

### 错误响应

| code | message | 说明 |
|------|---------|------|
| 20002 | Error occurred while binding the request body to the struct. : user xxx already exists | 用户已存在 |
| 10002 | Database error : ... | 数据库异常 |

---

## 2. 登录

```
POST /api/v1/auth/login
Content-Type: application/json
```

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名（邮箱） |
| password | string | 是 | 密码 |

### 请求示例

```json
{
  "username": "user@example.com",
  "password": "123456"
}
```

### 成功响应

```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### 后续请求

登录成功后，在需要鉴权的接口请求头中携带 token：

```
Authorization: eyJhbGciOiJIUzI1NiIs...
```
```

| 字段 | 类型 | 说明 |
|------|------|------|
| data.token | string | JWT token，有效期 7 天 |

### 错误响应

| code | message | 说明 |
|------|---------|------|
| 10002 | Database error : ... | 数据库异常 |
| 20101 | User not existed : user xxx not found | 用户不存在 |
| 20102 | The password was incorrect. | 密码错误 |
| 20005 | Error occurred while handling the auth token : ... | token 生成失败 |

---

## 其他认证相关接口

| 方法 | 路径 | 说明 | 鉴权 |
|------|------|------|------|
| POST | /api/v1/auth/set_role/:id | 设置用户角色 | 管理员 |
