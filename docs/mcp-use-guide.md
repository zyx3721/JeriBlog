## MCP (Model Context Protocol) 使用指南

### 1. MCP 是什么

**MCP (Model Context Protocol)** 是一个标准化协议，允许 AI 客户端（如 Claude Code、Cherry Studio）通过 HTTP 接口调用服务器端的工具（Tools），实现对博客系统的远程管理。

**核心特点**：
- 🔧 **9 个管理工具**：文章、分类/标签、评论、友链、RSS、统计、动态、用户、访问日志
- 🔐 **Bearer Token 认证**：使用 MCP Secret 进行身份验证
- 🌐 **HTTP 流式传输**：支持实时响应
- ⚠️ **包含高危操作**：删除文章、删除用户等不可逆操作

---

### 2. 配置步骤

#### 2.1 后端配置（获取 MCP Secret）

1. 登录管理后台
2. 进入 **设置 → AI 设置**
3. 找到 **MCP** 区域
4. 复制 **Secret**（系统自动生成）
5. 如需重置，点击 **重置** 按钮（⚠️ 仅超级管理员可操作）

**代码位置**：
- 前端界面：`admin/src/views/setting/components/AISettingsTab.vue` (lines 75-96)
- 后端认证：`server/api/middleware/auth.go` (lines 106-118)

#### 2.2 客户端配置

**方式一：Claude Code**

编辑 `%USERPROFILE%/.claude.json` 文件，添加：

```json
{
  "mcpServers": {
    "jeriblog": {
      "type": "http",
      "url": "https://你的后端地址/mcp",
      "headers": {
        "Authorization": "Bearer <你的MCP_Secret>"
      }
    }
  }
}
```

**方式二：Cherry Studio**

1. 打开设置
2. 找到 **MCP 服务器配置**
3. 添加新服务器：
   - **名称**：自定义（如 JeriBlog）
   - **类型**：可流式传输的 HTTP
   - **URL**：`https://你的后端地址/mcp`
   - **请求头**：`Authorization=Bearer <你的MCP_Secret>`

---

### 3. 技术架构

#### 3.1 路由注册

**代码位置**：`server/api/router/router.go` (line 135)

```go
// MCP 路由（需要 MCP Secret 认证）
mcpHandler := gin.WrapH(mcpserver.NewPublicHandler(
    articleService, categoryService, tagService, commentService,
    friendService, rssFeedService, momentService, userService, statsService,
))
r.Any("/mcp", middleware.MCPAuth(conf), mcpHandler)
```

#### 3.2 认证机制

**代码位置**：`server/api/middleware/auth.go` (lines 106-118)

```go
func MCPAuth(conf *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        token, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
        if ok && token != "" && conf != nil && conf.AI.MCPSecret != "" &&
            subtle.ConstantTimeCompare([]byte(token), []byte(conf.AI.MCPSecret)) == 1 {
            c.Next()
            return
        }
        response.Error(c, errcode.Unauthorized.WithDetails("MCP 认证失败"))
        c.Abort()
    }
}
```

**认证流程**：
1. 从 HTTP Header 提取 `Authorization: Bearer <token>`
2. 与数据库中的 `conf.AI.MCPSecret` 进行常量时间比较
3. 验证通过则放行，否则返回 401 错误

#### 3.3 工具注册

**代码位置**：`server/pkg/mcp/register.go` (lines 10-90)

系统注册了 9 个聚合工具：

```go
func (s *publicServer) registerTools(server *sdkmcp.Server) {
    // 1. article_manage - 文章管理
    sdkmcp.AddTool(server, &sdkmcp.Tool{
        Name:        "article_manage",
        Description: "文章管理。action：list/get/create/update/delete。",
        InputSchema: tools.ArticleManageInputSchema(),
    }, articleWrapper.ManageArticle)

    // 2. taxonomy_manage - 分类/标签管理
    // 3. comment_manage - 评论管理
    // 4. friend_manage - 友链管理
    // 5. rssfeed_manage - RSS订阅管理
    // 6. stats_query - 统计查询
    // 7. moment_manage - 动态管理
    // 8. user_manage - 用户管理
    // 9. visit_manage - 访问日志管理
}
```

---

### 4. 工具详细说明

#### 4.1 文章管理 (article_manage)

**支持操作**：

| Action   | 说明         | 必需参数 | 可选参数 | 风险 |
| -------- | ------------ | ------- | ------- | ---- |
| `list`   | 获取文章列表 | - | `page`, `page_size`, `keyword`, `category_id`, `tag_id`, `is_publish` | 无   |
| `get`    | 获取文章详情 | `id` | - | 无   |
| `create` | 创建文章 | `title`, `content` | `summary`, `category_id`, `tag_ids`, `is_publish`, `cover_url` | 低   |
| `update` | 更新文章     | `id` | `title`, `content`, `summary`, `category_id`, `tag_ids`, `is_publish`, `cover_url` | 中   |
| `delete` | 删除文章     | `id` | - | ⚠️ 高 |

**参数说明**：
- `page`: 页码（默认 1）
- `page_size`: 每页数量（默认 20）
- `keyword`: 搜索关键词（标题/内容）
- `category_id`: 分类 ID
- `tag_id`: 标签 ID
- `is_publish`: 发布状态（true/false）
- `title`: 文章标题
- `content`: 文章内容（Markdown 格式）
- `summary`: 文章摘要
- `tag_ids`: 标签 ID 数组
- `cover_url`: 封面图片 URL

**代码位置**：`server/pkg/mcp/tools/article.go`

#### 4.2 分类/标签管理 (taxonomy_manage)

**支持操作**：

| Target     | Action          | 必需参数 | 可选参数 | 说明             | 风险 |
| ---------- | --------------- | ------- | ------- | ---------------- | ---- |
| `category` | `list`          | - | `page`, `page_size` | 获取分类列表     | 无   |
| `category` | `get`           | `id` | - | 获取分类详情     | 无   |
| `category` | `create`        | `name`, `slug` | `description`, `icon` | 创建分类         | 低   |
| `category` | `update`        | `id` | `name`, `slug`, `description`, `icon` | 更新分类         | 中   |
| `category` | `delete`        | `id` | - | 删除分类         | ⚠️ 高 |
| `category` | `list_articles` | `id` | `page`, `page_size` | 获取分类下的文章 | 无   |
| `tag`      | `list`          | - | `page`, `page_size` | 获取标签列表     | 无   |
| `tag`      | `get`           | `id` | - | 获取标签详情     | 无   |
| `tag`      | `create`        | `name`, `slug` | `description`, `color` | 创建标签         | 低   |
| `tag`      | `update`        | `id` | `name`, `slug`, `description`, `color` | 更新标签         | 中   |
| `tag`      | `delete`        | `id` | - | 删除标签         | ⚠️ 高 |
| `tag`      | `list_articles` | `id` | `page`, `page_size` | 获取标签下的文章 | 无   |

**参数说明**：
- `name`: 分类/标签名称
- `slug`: URL 友好的标识符（如 "golang"）
- `description`: 描述信息
- `icon`: 分类图标（仅分类）
- `color`: 标签颜色（仅标签，如 "#409EFF"）

**代码位置**：`server/pkg/mcp/tools/taxonomy.go`

#### 4.3 评论管理 (comment_manage)

**支持操作**：

| Action          | 必需参数 | 可选参数 | 说明              | 风险 |
| --------------- | ------- | ------- | ----------------- | ---- |
| `list`          | - | `page`, `page_size`, `article_id`, `status` | 获取评论列表      | 无   |
| `get`           | `id` | - | 获取评论详情      | 无   |
| `toggle_status` | `id` | - | 切换显示/隐藏状态 | 低   |
| `restore`       | `id` | - | 恢复已删除的评论  | 低   |
| `delete`        | `id` | - | 删除评论（硬删除）| ⚠️ 高 |

**参数说明**：
- `article_id`: 文章 ID（筛选指定文章的评论）
- `status`: 评论状态（`visible`/`hidden`）

**代码位置**：`server/pkg/mcp/tools/comment.go`

#### 4.4 友链管理 (friend_manage)

**支持操作**：

| Action   | 必需参数 | 可选参数 | 说明         | 风险 |
| -------- | ------- | ------- | ------------ | ---- |
| `list`   | - | `page`, `page_size`, `type_id`, `status` | 获取友链列表 | 无   |
| `get`    | `id` | - | 获取友链详情 | 无   |
| `create` | `name`, `url`, `avatar` | `description`, `type_id`, `status` | 创建友链     | 低   |
| `update` | `id` | `name`, `url`, `avatar`, `description`, `type_id`, `status` | 更新友链     | 中   |
| `delete` | `id` | - | 删除友链     | ⚠️ 高 |

**参数说明**：
- `name`: 友链名称
- `url`: 友链地址
- `avatar`: 头像 URL
- `description`: 友链描述
- `type_id`: 友链类型 ID
- `status`: 友链状态（`pending`/`approved`/`rejected`）

**代码位置**：`server/pkg/mcp/tools/friend.go`

#### 4.5 RSS订阅管理 (rssfeed_manage)

**支持操作**：

| Action          | 必需参数 | 可选参数 | 说明             | 风险 |
| --------------- | ------- | ------- | ---------------- | ---- |
| `list`          | - | `page`, `page_size`, `is_read`, `feed_url` | 获取RSS文章列表  | 无   |
| `mark_read`     | `id` | - | 标记单篇文章已读 | 低   |
| `mark_all_read` | - | - | 标记所有文章已读 | 中   |

**参数说明**：
- `is_read`: 已读状态（true/false）
- `feed_url`: RSS 源地址（筛选指定源的文章）

**代码位置**：`server/pkg/mcp/tools/rssfeed.go`

#### 4.6 统计查询 (stats_query)

**支持操作**：

| Action      | 必需参数 | 可选参数 | 说明               | 风险 |
| ----------- | ------- | ------- | ------------------ | ---- |
| `dashboard` | - | - | 查询仪表盘概览数据 | 无   |
| `trend`     | - | `days` | 查询访问趋势数据   | 无   |

**参数说明**：
- `days`: 查询天数（默认 7 天）

**返回数据**：
- `dashboard`: 文章数、评论数、访问量、用户数等
- `trend`: 每日访问量、PV、UV 趋势图数据

**特点**：只读工具，无风险

**代码位置**：`server/pkg/mcp/tools/stats.go`

#### 4.7 动态管理 (moment_manage)

**支持操作**：

| Action   | 必需参数 | 可选参数 | 说明         | 风险 |
| -------- | ------- | ------- | ------------ | ---- |
| `list`   | - | `page`, `page_size`, `is_publish` | 获取动态列表 | 无   |
| `get`    | `id` | - | 获取动态详情 | 无   |
| `create` | `content` | `is_publish`, `images` | 创建动态     | 低   |
| `update` | `id` | `content`, `is_publish`, `images` | 更新动态     | 中   |
| `delete` | `id` | - | 删除动态     | ⚠️ 高 |

**参数说明**：
- `content`: 动态内容（支持 Markdown）
- `is_publish`: 发布状态（true/false）
- `images`: 图片 URL 数组

**代码位置**：`server/pkg/mcp/tools/moment.go`

#### 4.8 用户管理 (user_manage)

**支持操作**：

| Action   | 必需参数 | 可选参数 | 说明         | 风险 |
| -------- | ------- | ------- | ------------ | ---- |
| `list`   | - | `page`, `page_size`, `role`, `is_enabled` | 获取用户列表 | 无   |
| `get`    | `id` | - | 获取用户详情 | 无   |
| `create` | `nickname`, `email`, `password` | `role`, `avatar`, `website` | 创建用户     | 中   |
| `update` | `id` | `nickname`, `email`, `role`, `is_enabled`, `avatar`, `website` | 更新用户     | 中   |
| `delete` | `id` | - | 删除用户     | ⚠️ 高 |

**参数说明**：
- `nickname`: 用户昵称
- `email`: 邮箱地址
- `password`: 密码（仅创建时需要）
- `role`: 用户角色（`super_admin`/`admin`/`user`/`guest`）
- `is_enabled`: 启用状态（true/false）
- `avatar`: 头像 URL
- `website`: 个人网站

**代码位置**：`server/pkg/mcp/tools/user.go`

---

#### 4.9 访问日志管理 (visit_manage)

**支持操作**：

| Action         | 必需参数 | 可选参数 | 说明             | 风险 |
| -------------- | ------- | ------- | ---------------- | ---- |
| `list`         | - | `page`, `page_size`, `keyword`, `visitor_id`, `ip`, `exclude_ips`, `location`, `browser`, `os`, `start_time`, `end_time` | 获取访问日志列表 | 无   |
| `delete`       | `id` | - | 删除单条访问日志 | ⚠️ 高 |
| `batch_delete` | `ids` | - | 批量删除访问日志 | ⚠️ 高 |

**参数说明**：
- `page`: 页码（默认 1）
- `page_size`: 每页数量（默认 20）
- `keyword`: 搜索关键词（页面 URL）
- `visitor_id`: 访客唯一标识
- `ip`: IP 地址
- `exclude_ips`: 排除的 IP 地址，多个用逗号分隔
- `location`: 地理位置
- `browser`: 浏览器
- `os`: 操作系统
- `start_time`: 开始时间（格式：2006-01-02）
- `end_time`: 结束时间（格式：2006-01-02）
- `ids`: 要删除的访问日志 ID 列表（数组）

**代码位置**：`server/pkg/mcp/tools/visit.go`

---

### 5. 安全注意事项

#### 5.1 高危操作

以下操作不可逆，使用前请确认：

- ⚠️ **删除文章** (`article_manage` → `delete`)
- ⚠️ **删除分类/标签** (`taxonomy_manage` → `delete`)
- ⚠️ **删除评论** (`comment_manage` → `delete`)
- ⚠️ **删除友链** (`friend_manage` → `delete`)
- ⚠️ **删除动态** (`moment_manage` → `delete`)
- ⚠️ **删除用户** (`user_manage` → `delete`)
- ⚠️ **删除访问日志** (`visit_manage` → `delete`/`batch_delete`)

#### 5.2 Secret 管理

- ✅ **定期重置**：建议定期更换 MCP Secret
- ✅ **权限控制**：仅超级管理员可重置 Secret
- ✅ **安全存储**：不要将 Secret 提交到版本控制
- ✅ **HTTPS 传输**：生产环境必须使用 HTTPS

#### 5.3 模型兼容性

根据文档测试结果：

| 模型          | 兼容性 | 说明             |
| ------------- | ------ | ---------------- |
| GPT-5.4       | ✅ 优秀 | 极低失误率       |
| DeepSeek-V3.2 | ✅ 优秀 | 极低失误率       |
| MiniMax-M2.5  | ✅ 优秀 | 极低失误率       |
| Kimi-K2.5     | ⚠️ 一般 | 需多次调用才正确 |

---

### 6. 使用示例

#### 6.1 文章管理示例

##### 示例 1：查询文章列表

```
用户：帮我查看最近的 5 篇文章

Claude：我来调用 article_manage 工具查询...
[调用 MCP Tool: article_manage]
{
  "action": "list",
  "payload": {
    "page": 1,
    "page_size": 5
  }
}

返回结果：
1. MCP 使用指南 (ID: 123, 分类: 技术, 发布时间: 2026-05-01)
2. Go 语言实战 (ID: 122, 分类: 编程, 发布时间: 2026-04-30)
3. Docker 容器化部署 (ID: 121, 分类: 运维, 发布时间: 2026-04-29)
4. Vue3 组件开发 (ID: 120, 分类: 前端, 发布时间: 2026-04-28)
5. MySQL 性能优化 (ID: 119, 分类: 数据库, 发布时间: 2026-04-27)
```

##### 示例 2：获取文章详情

```
用户：查看 ID 为 123 的文章详情

Claude：我来获取文章详情...
[调用 MCP Tool: article_manage]
{
  "action": "get",
  "payload": {
    "id": 123
  }
}

返回结果：
标题：MCP 使用指南
分类：技术
标签：MCP, Claude, AI
发布状态：已发布
创建时间：2026-05-01 10:30:00
内容：（文章内容）...
```

##### 示例 3：创建草稿文章

```
用户：帮我创建一篇标题为"Redis 缓存实战"的草稿文章

Claude：我来创建草稿文章...
[调用 MCP Tool: article_manage]
{
  "action": "create",
  "payload": {
    "title": "Redis 缓存实战",
    "content": "",
    "is_publish": false
  }
}

创建成功！文章 ID: 124
```

##### 示例 4：更新文章并发布

```
用户：将 ID 为 124 的文章内容更新为"Redis 是一个高性能的键值数据库..."，并发布

Claude：我来更新文章并发布...
[调用 MCP Tool: article_manage]
{
  "action": "update",
  "payload": {
    "id": 124,
    "content": "Redis 是一个高性能的键值数据库...",
    "is_publish": true
  }
}

更新成功！文章已发布。
```

##### 示例 5：删除文章

```
用户：删除 ID 为 124 的文章

Claude：⚠️ 警告：删除操作不可逆，确认要删除吗？

用户：确认删除

Claude：我来删除文章...
[调用 MCP Tool: article_manage]
{
  "action": "delete",
  "payload": {
    "id": 124
  }
}

删除成功！
```

#### 6.2 分类/标签管理示例

##### 示例 6：获取所有分类

```
用户：查看所有分类

Claude：我来获取分类列表...
[调用 MCP Tool: taxonomy_manage]
{
  "target": "category",
  "action": "list",
  "payload": {}
}

返回结果：
1. 技术 (ID: 1, slug: tech, 文章数: 25)
2. 编程 (ID: 2, slug: programming, 文章数: 18)
3. 运维 (ID: 3, slug: devops, 文章数: 12)
4. 前端 (ID: 4, slug: frontend, 文章数: 15)
```

##### 示例 7：创建新分类

```
用户：创建一个名为"数据库"的分类，slug 为 "database"

Claude：我来创建分类...
[调用 MCP Tool: taxonomy_manage]
{
  "target": "category",
  "action": "create",
  "payload": {
    "name": "数据库",
    "slug": "database",
    "description": "数据库相关技术文章"
  }
}

创建成功！分类 ID: 5
```

##### 示例 8：获取分类下的文章

```
用户：查看"技术"分类下的所有文章

Claude：我来查询分类下的文章...
[调用 MCP Tool: taxonomy_manage]
{
  "target": "category",
  "action": "list_articles",
  "payload": {
    "id": 1,
    "page": 1,
    "page_size": 10
  }
}

返回结果：
技术分类下共有 25 篇文章：
1. MCP 使用指南 (ID: 123)
2. API 设计规范 (ID: 118)
...
```

##### 示例 9：创建标签

```
用户：创建一个名为"Redis"的标签，颜色为红色

Claude：我来创建标签...
[调用 MCP Tool: taxonomy_manage]
{
  "target": "tag",
  "action": "create",
  "payload": {
    "name": "Redis",
    "slug": "redis",
    "color": "#F56C6C"
  }
}

创建成功！标签 ID: 15
```

##### 示例 10：更新标签

```
用户：将 ID 为 15 的标签颜色改为蓝色

Claude：我来更新标签...
[调用 MCP Tool: taxonomy_manage]
{
  "target": "tag",
  "action": "update",
  "payload": {
    "id": 15,
    "color": "#409EFF"
  }
}

更新成功！
```

#### 6.3 评论管理示例

##### 示例 11：获取评论列表

```
用户：查看最近的 10 条评论

Claude：我来获取评论列表...
[调用 MCP Tool: comment_manage]
{
  "action": "list",
  "payload": {
    "page": 1,
    "page_size": 10
  }
}

返回结果：
1. 用户 Alice：这篇文章写得很好！(文章: MCP 使用指南, 时间: 2026-05-03 10:30)
2. 用户 Bob：感谢分享 (文章: Go 语言实战, 时间: 2026-05-03 09:15)
...
```

##### 示例 12：获取指定文章的评论

```
用户：查看文章 ID 为 123 的所有评论

Claude：我来查询该文章的评论...
[调用 MCP Tool: comment_manage]
{
  "action": "list",
  "payload": {
    "article_id": 123,
    "page": 1,
    "page_size": 20
  }
}

返回结果：
文章"MCP 使用指南"共有 5 条评论：
1. Alice：这篇文章写得很好！
2. Bob：请问如何配置 MCP？
...
```

##### 示例 13：隐藏评论

```
用户：隐藏 ID 为 50 的评论

Claude：我来切换评论状态...
[调用 MCP Tool: comment_manage]
{
  "action": "toggle_status",
  "payload": {
    "id": 50
  }
}

操作成功！评论已隐藏。
```

##### 示例 14：恢复已删除的评论

```
用户：恢复 ID 为 48 的评论

Claude：我来恢复评论...
[调用 MCP Tool: comment_manage]
{
  "action": "restore",
  "payload": {
    "id": 48
  }
}

恢复成功！
```

##### 示例 15：删除评论

```
用户：删除 ID 为 50 的评论

Claude：⚠️ 警告：删除操作不可逆，确认要删除吗？

用户：确认删除

Claude：我来删除评论...
[调用 MCP Tool: comment_manage]
{
  "action": "delete",
  "payload": {
    "id": 50
  }
}

删除成功！
```

#### 6.4 友链管理示例

##### 示例 16：获取友链列表

```
用户：查看所有友链

Claude：我来获取友链列表...
[调用 MCP Tool: friend_manage]
{
  "action": "list",
  "payload": {}
}

返回结果：
1. 张三的博客 (URL: https://zhangsan.com, 状态: 已通过)
2. 李四技术站 (URL: https://lisi.tech, 状态: 待审核)
3. 王五的小站 (URL: https://wangwu.cn, 状态: 已通过)
```

##### 示例 17：创建友链

```
用户：添加一个友链，名称"小明的博客"，地址 https://xiaoming.com

Claude：我来创建友链...
[调用 MCP Tool: friend_manage]
{
  "action": "create",
  "payload": {
    "name": "小明的博客",
    "url": "https://xiaoming.com",
    "avatar": "https://xiaoming.com/avatar.jpg",
    "description": "专注前端技术分享",
    "status": "approved"
  }
}

创建成功！友链 ID: 10
```

##### 示例 18：更新友链状态

```
用户：将 ID 为 2 的友链状态改为已通过

Claude：我来更新友链状态...
[调用 MCP Tool: friend_manage]
{
  "action": "update",
  "payload": {
    "id": 2,
    "status": "approved"
  }
}

更新成功！
```

##### 示例 19：删除友链

```
用户：删除 ID 为 10 的友链

Claude：⚠️ 警告：删除操作不可逆，确认要删除吗？

用户：确认删除

Claude：我来删除友链...
[调用 MCP Tool: friend_manage]
{
  "action": "delete",
  "payload": {
    "id": 10
  }
}

删除成功！
```

#### 6.5 RSS订阅管理示例

##### 示例 20：获取未读 RSS 文章

```
用户：查看所有未读的 RSS 文章

Claude：我来获取未读文章列表...
[调用 MCP Tool: rssfeed_manage]
{
  "action": "list",
  "payload": {
    "is_read": false,
    "page": 1,
    "page_size": 20
  }
}

返回结果：
共有 15 篇未读文章：
1. Go 1.22 新特性解析 (来源: Go 官方博客, 时间: 2026-05-03)
2. Docker 最佳实践 (来源: Docker Blog, 时间: 2026-05-02)
...
```

##### 示例 21：标记单篇文章已读

```
用户：将 ID 为 100 的 RSS 文章标记为已读

Claude：我来标记文章已读...
[调用 MCP Tool: rssfeed_manage]
{
  "action": "mark_read",
  "payload": {
    "id": 100
  }
}

标记成功！
```

##### 示例 22：标记所有文章已读

```
用户：将所有 RSS 文章标记为已读

Claude：我来标记所有文章已读...
[调用 MCP Tool: rssfeed_manage]
{
  "action": "mark_all_read",
  "payload": {}
}

操作成功！所有文章已标记为已读。
```

#### 6.6 统计查询示例

##### 示例 23：查询仪表盘数据

```
用户：查看网站统计数据

Claude：我来获取仪表盘数据...
[调用 MCP Tool: stats_query]
{
  "action": "dashboard",
  "payload": {}
}

返回结果：
📊 网站统计概览
- 文章总数：125 篇
- 评论总数：1,234 条
- 用户总数：567 人
- 今日访问：1,890 次
- 总访问量：125,678 次
```

##### 示例 24：查询访问趋势

```
用户：查看最近 7 天的访问趋势

Claude：我来获取访问趋势数据...
[调用 MCP Tool: stats_query]
{
  "action": "trend",
  "payload": {
    "days": 7
  }
}

返回结果：
📈 最近 7 天访问趋势
- 2026-05-03: 1,890 次 (PV: 2,345, UV: 1,234)
- 2026-05-02: 1,756 次 (PV: 2,123, UV: 1,156)
- 2026-05-01: 1,623 次 (PV: 1,987, UV: 1,089)
...
```

#### 6.7 动态管理示例

##### 示例 25：获取动态列表

```
用户：查看最近的 5 条动态

Claude：我来获取动态列表...
[调用 MCP Tool: moment_manage]
{
  "action": "list",
  "payload": {
    "page": 1,
    "page_size": 5
  }
}

返回结果：
1. 今天天气真好！☀️ (时间: 2026-05-03 10:30)
2. 刚完成了一个新项目 🎉 (时间: 2026-05-02 15:20)
3. 分享一张美图 📷 (时间: 2026-05-01 09:15)
...
```

##### 示例 26：创建动态

```
用户：发布一条动态"今天学习了 MCP 协议"

Claude：我来创建动态...
[调用 MCP Tool: moment_manage]
{
  "action": "create",
  "payload": {
    "content": "今天学习了 MCP 协议",
    "is_publish": true
  }
}

创建成功！动态 ID: 50
```

##### 示例 27：更新动态

```
用户：将 ID 为 50 的动态内容改为"今天深入学习了 MCP 协议，收获很多！"

Claude：我来更新动态...
[调用 MCP Tool: moment_manage]
{
  "action": "update",
  "payload": {
    "id": 50,
    "content": "今天深入学习了 MCP 协议，收获很多！"
  }
}

更新成功！
```

##### 示例 28：删除动态

```
用户：删除 ID 为 50 的动态

Claude：⚠️ 警告：删除操作不可逆，确认要删除吗？

用户：确认删除

Claude：我来删除动态...
[调用 MCP Tool: moment_manage]
{
  "action": "delete",
  "payload": {
    "id": 50
  }
}

删除成功！
```

#### 6.8 用户管理示例

##### 示例 29：获取用户列表

```
用户：查看所有管理员用户

Claude：我来获取管理员列表...
[调用 MCP Tool: user_manage]
{
  "action": "list",
  "payload": {
    "role": "admin",
    "page": 1,
    "page_size": 20
  }
}

返回结果：
共有 3 位管理员：
1. 张三 (邮箱: zhangsan@example.com, 角色: 超级管理员, 状态: 启用)
2. 李四 (邮箱: lisi@example.com, 角色: 管理员, 状态: 启用)
3. 王五 (邮箱: wangwu@example.com, 角色: 管理员, 状态: 禁用)
```

##### 示例 30：获取用户详情

```
用户：查看 ID 为 5 的用户详情

Claude：我来获取用户详情...
[调用 MCP Tool: user_manage]
{
  "action": "get",
  "payload": {
    "id": 5
  }
}

返回结果：
用户信息：
- 昵称：张三
- 邮箱：zhangsan@example.com
- 角色：超级管理员
- 状态：启用
- 注册时间：2026-01-01 10:00:00
- 最后登录：2026-05-03 09:30:00
```

##### 示例 31：创建用户

```
用户：创建一个新用户，昵称"小红"，邮箱 xiaohong@example.com，密码 123456

Claude：我来创建用户...
[调用 MCP Tool: user_manage]
{
  "action": "create",
  "payload": {
    "nickname": "小红",
    "email": "xiaohong@example.com",
    "password": "123456",
    "role": "user"
  }
}

创建成功！用户 ID: 20
```

##### 示例 32：更新用户角色

```
用户：将 ID 为 20 的用户角色改为管理员

Claude：我来更新用户角色...
[调用 MCP Tool: user_manage]
{
  "action": "update",
  "payload": {
    "id": 20,
    "role": "admin"
  }
}

更新成功！
```

##### 示例 33：禁用用户

```
用户：禁用 ID 为 20 的用户

Claude：我来禁用用户...
[调用 MCP Tool: user_manage]
{
  "action": "update",
  "payload": {
    "id": 20,
    "is_enabled": false
  }
}

操作成功！用户已禁用。
```

##### 示例 34：删除用户

```
用户：删除 ID 为 20 的用户

Claude：⚠️ 警告：删除操作不可逆，确认要删除吗？

用户：确认删除

Claude：我来删除用户...
[调用 MCP Tool: user_manage]
{
  "action": "delete",
  "payload": {
    "id": 20
  }
}

删除成功！
```

---

#### 6.9 访问日志管理示例

##### 示例 35：列出访问日志

```
用户：列出最近的访问日志

Claude：我来获取访问日志列表...
[调用 MCP Tool: visit_manage]
{
  "action": "list",
  "payload": {
    "page": 1,
    "page_size": 20
  }
}

已获取 20 条访问日志：
1. ID: 1001 - IP: 192.168.1.100 - 页面: /article/123 - 时间: 2026-05-03 10:30:45
2. ID: 1002 - IP: 192.168.1.101 - 页面: /article/124 - 时间: 2026-05-03 10:31:20
...
```

##### 示例 36：筛选特定 IP 的访问日志

```
用户：查看 IP 为 192.168.1.100 的访问记录

Claude：我来筛选该 IP 的访问日志...
[调用 MCP Tool: visit_manage]
{
  "action": "list",
  "payload": {
    "ip": "192.168.1.100",
    "page": 1,
    "page_size": 20
  }
}

已找到 5 条该 IP 的访问记录：
1. ID: 1001 - 页面: /article/123 - 浏览器: Chrome - 时间: 2026-05-03 10:30:45
2. ID: 1005 - 页面: /article/125 - 浏览器: Chrome - 时间: 2026-05-03 11:15:30
...
```

##### 示例 37：按时间范围筛选访问日志

```
用户：查看 2026-05-01 到 2026-05-03 的访问日志

Claude：我来获取该时间段的访问日志...
[调用 MCP Tool: visit_manage]
{
  "action": "list",
  "payload": {
    "start_time": "2026-05-01",
    "end_time": "2026-05-03",
    "page": 1,
    "page_size": 50
  }
}

已获取该时间段的 50 条访问日志（共 150 条）。
```

##### 示例 38：删除单条访问日志

```
用户：删除 ID 为 1001 的访问日志

Claude：⚠️ 警告：删除操作不可逆，确认要删除吗？

用户：确认删除

Claude：我来删除访问日志...
[调用 MCP Tool: visit_manage]
{
  "action": "delete",
  "payload": {
    "id": 1001
  }
}

删除成功！
```

##### 示例 39：批量删除访问日志

```
用户：批量删除 ID 为 1001、1002、1003 的访问日志

Claude：⚠️ 警告：批量删除操作不可逆，确认要删除这 3 条记录吗？

用户：确认删除

Claude：我来批量删除访问日志...
[调用 MCP Tool: visit_manage]
{
  "action": "batch_delete",
  "payload": {
    "ids": [1001, 1002, 1003]
  }
}

批量删除成功！已删除 3 条访问日志。
```

---

### 7. 配置文件位置

| 配置项       | 位置                                                   | 说明                             |
| ------------ | ------------------------------------------------------ | -------------------------------- |
| MCP Secret   | 数据库 `settings` 表                                   | `group='ai'`, `key='mcp_secret'` |
| 前端配置界面 | `admin/src/views/setting/components/AISettingsTab.vue` | AI 设置页面                      |
| 后端路由     | `server/api/router/router.go`                          | `/mcp` 路由                      |
| 认证中间件   | `server/api/middleware/auth.go`                        | `MCPAuth()`                      |
| 工具注册     | `server/pkg/mcp/register.go`                           | 9 个工具定义                     |
| 工具实现     | `server/pkg/mcp/tools/*.go`                            | 各工具具体实现                   |

---

### 8. 常见问题

**Q1: MCP Secret 在哪里查看？**
A: 管理后台 → 设置 → AI 设置 → MCP 区域

**Q2: 如何重置 MCP Secret？**
A: 仅超级管理员可在 AI 设置页面点击"重置"按钮

**Q3: 为什么连接失败？**
A: 检查以下项：
- URL 是否正确（需包含 `/mcp` 路径）
- Secret 是否正确复制
- 后端服务是否正常运行
- 是否使用 HTTPS（生产环境）

**Q4: 哪些操作有风险？**
A: 所有 `delete` 操作均不可逆，使用前请确认

**Q5: 可以同时配置多个客户端吗？**
A: 可以，所有客户端使用同一个 MCP Secret

---

## 总结

MCP 是 JeriBlog 提供的强大远程管理接口，允许通过 AI 客户端（如 Claude Code）直接管理博客内容。核心特点：

- ✅ **9 个管理工具**：覆盖文章、分类、评论、友链、RSS、统计、动态、用户、访问日志
- ✅ **Bearer Token 认证**：使用 MCP Secret 保证安全
- ✅ **HTTP 流式传输**：支持实时响应
- ⚠️ **包含高危操作**：删除操作不可逆，需谨慎使用

**推荐使用场景**：
- 批量管理文章
- 快速查询统计数据
- 远程内容审核
- 自动化内容发布
- 访问日志分析与清理