# a-easy-community
community like zhihu

## 基础信息
- **Base URL**: `http://localhost:8080` (默认本地运行端口，根据部署情况可能变化)

## 技术栈

- **语言**: Go (1.25)
- **Web 框架**: Gin 
- **数据库**: MySQL 8.0
- **缓存**: Redi
- **ORM**: GORM
- **配置**: Viper
- **日志**: Zap + Lumberjack (归档)
- **鉴权**: JWT + Bcrypt
- **并发控制**: Singleflight + 缓存空值 + 随机过期时间 
- **ai**: Glm4.7 flash
- **UUID** : 生成唯一标识符
- **HTML过滤**：防止XSS攻击，用于净化用户输入
- **防csrf**：`c.SetSameSite(http.SameSiteStrictMode)`确保refreshcookie不会在其他网站调用
- **websocket**：（目前还未实现）

## 认证方式

- **Access Token**: 除登录和注册外的接口需要 JWT 认证。请在 Header 中携带：
  
  ```
  Authorization: Bearer <your_token>
  ```
- **Refresh Token**: 登录成功后会通过 Cookie 下发 `refresh_token`，用于刷新 Access Token。

---

## 1. 账户与认证 (Account)

### 注册
- **URL**: `/account/register`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "account": "userAccount",
    "password": "password",
    "name": "昵称"
  }
  ```

### 登录
- **URL**: `/account/login`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "account": "userAccount",
    "password": "password"
  }
  ```
- **Response**:
  ```json
  {
    "code": 200,
    "data": { "token": "eyJhbGciOiJIUzI1..." },
    "msg": "success"
  }
  ```
  *(注：Refresh Token 会写入 Cookie)*

### 刷新 Token
- **URL**: `/account/refresh`
- **Method**: `POST`
- **Cookie**: 需要携带 `refresh_token`

---

*以下接口均需携带 Bearer Token*

## 2. 用户个人中心 (User Profile)

| 接口功能         | URL                                  | Method   | 说明                                      |
| :--------------- | :----------------------------------- | :------- | :---------------------------------------- |
| **获取个人信息** | `/account/protected/profile`         | `GET`    | 获取当前登录用户信息                      |
| **修改用户名**   | `/account/protected/username`        | `PATCH`  |                                           |
| **修改头像**     | `/account/protected/avatar`          | `POST`   | 上传头像文件                              |
| **修改简介**     | `/account/protected/introduction`    | `PATCH`  |                                           |
| **修改密码**     | `/account/protected/password-change` | `POST`   | **限流**: 5秒/1次                         |
| **退出登录**     | `/account/protected/logout`          | `POST`   | token和refresh token会存入redis，等待过期 |
| **注销账户**     | `/account/protected/delete-user`     | `DELETE` | 软删除                                    |

---

## 3. 社交与用户管理 (Social & Management)
| 接口功能         | URL                             | Method | 说明               |
| :--------------- | :------------------------------ | :----- | :----------------- |
| **查看他人主页** | `/account/protected/users/:Id`  | `GET`  | `:Id` 为用户ID     |
| **关注/取关**    | `/account/protected/follow/:Id` | `POST` | `:Id` 为目标用户ID |
| **禁言用户**     | `/account/protected/muted/:Id`  | `POST` | 管理员功能         |
| **设置VIP**      | `/account/protected/vip/:Id`    | `POST` | 管理员功能         |

---

## 4. 帖子与内容 (Posts)
### 获取帖子列表
- **URL**: `/account/protected/posts`
- **Method**: `GET`
- **Query Params**: `page` (默认 1)

### 发布帖子
- **URL**: `/account/protected/posts`
- **Method**: `POST`
- **限流**: 5秒/1次
- **Body**:
  ```json
  {
    "title": "帖子标题",
    "content": "帖子内容"
  }
  ```

### 帖子详情
- **URL**: `/account/protected/posts/:postId`
- **Method**: `GET`

### 搜索帖子
- **URL**: `/account/protected/search`
- **Method**: `GET`
- **Query Params**: `keyword` (假设参数)

### 帖子交互
| 接口功能        | URL                                        | Method   | 限流    | 说明        |
| :-------------- | :----------------------------------------- | :------- | :------ | :---------- |
| **删除帖子**    | `/account/protected/posts/:postId`         | `DELETE` | -       |             |
| **点赞**        | `/account/protected/posts/:postId/like`    | `POST`   | 5秒/2次 |             |
| **AI总结(VIP)** | `/account/protected/posts/:postId/summary` | `POST`   | -       | VIP专属功能 |
| **设置付费贴**  | `/account/protected/paid-post/:postId`     | `POST`   | -       | 管理员设置  |

### 图片上传
- **URL**: `/account/protected/upload`
- **Method**: `POST`
- **Content-Type**: `multipart/form-data`
- **Form Param**: `file` (图片文件)
- **Response**: 返回图片访问 URL

---

## 5. 评论系统 (Comments)
### 发表评论
- **URL**: `/account/protected/posts/:postId`
- **Method**: `POST`
- **限流**: 3秒/1次
- **Body**:
  ```json
  {
    "content": "评论内容"
  }
  ```

### 删除评论
- **URL**: `/account/protected/posts/:postId/:posterId/:commentId`
- **Method**: `DELETE`

---

## 6. 其他 (Misc)

### 热度榜单
- **URL**: `/account/protected/hot_rank`
- **Method**: `GET`
- **说明**: 后台定时任务每5小时刷新一次热度。

### 静态资源

- **URL**: `/static/*`
- **Method**: `GET`
- **说明**: 访问上传的图片等静态文件。

### docker部署

- 借用了ai之力



(该文档在ai的基础上修改)