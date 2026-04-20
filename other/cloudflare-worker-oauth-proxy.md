# Cloudflare Worker OAuth 代理配置指南

## 一、为什么需要 OAuth 代理

由于国内网络环境限制，直接访问 GitHub、Google、Microsoft 的 OAuth API 可能会失败或超时。因此需要通过 Cloudflare Worker 搭建一个代理服务来中转这些请求。

## 二、代理服务的作用

Cloudflare Worker 代理会将以下请求转发到对应的 OAuth 服务：

| 代理路径 | 目标地址 | 用途 |
|---------|---------|------|
| `/github-api/*` | `https://api.github.com/*` | GitHub API 请求 |
| `/github/*` | `https://github.com/*` | GitHub OAuth 认证 |
| `/google-oauth2/*` | `https://oauth2.googleapis.com/*` | Google OAuth2 |
| `/google-api/*` | `https://www.googleapis.com/*` | Google API |
| `/google/*` | `https://accounts.google.com/*` | Google 账号认证 |
| `/microsoft-graph/*` | `https://graph.microsoft.com/*` | Microsoft Graph API |
| `/microsoft/*` | `https://login.microsoftonline.com/*` | Microsoft 登录 |

## 三、部署 Cloudflare Worker

### 3.1 注册 Cloudflare 账号

1. 访问：https://dash.cloudflare.com/sign-up
2. 使用邮箱注册（免费账号即可）
3. 验证邮箱并登录

### 3.2 创建 Worker

#### 步骤 1：进入 Workers 管理页面

1. 登录 Cloudflare 控制台后，点击左侧菜单 **"Workers & Pages"**
2. 点击右上角的 **"Create Application"** 按钮

#### 步骤 2：选择创建方式

在弹出的页面中，会看到两个选项卡：
- **Create Worker**（创建 Worker）
- **Create Pages**（创建 Pages）

**点击 "Create Worker" 选项卡**

#### 步骤 3：选择模板

在 "Create Worker" 页面中，会看到多个模板选项：
- **Start with Hello World!**（从 Hello World 开始）
- 其他预设模板...

**选择 "Start with Hello World!"**

然后点击右下角的 **"Create Worker"** 按钮。

#### 步骤 4：配置 Worker 名称

1. 在 "Name your Worker" 输入框中输入 Worker 名称
   - 建议命名：`jeriblog-oauth-proxy`
   - 命名规则：只能包含小写字母、数字和连字符
2. 点击右下角的蓝色 **"Deploy"** 按钮

**注意**：此时会部署一个默认的 Hello World 模板，这是正常的，我们稍后会替换代码。

#### 步骤 5：等待部署完成

部署成功后，会看到：
- 成功提示信息
- Worker 的访问地址（格式：`https://<worker-name>.<your-username>.workers.dev`）
- 页面中间会显示 Worker 的预览界面

**重要**：此时不要点击页面上的 "Edit code" 按钮（该按钮可能是灰色不可点击的）。

### 3.3 编辑 Worker 代码

#### 步骤 1：进入 Worker 管理页面

1. 部署成功后，点击页面左上角的 **"< Workers & Pages"** 返回列表
2. 或者直接点击左侧菜单的 **"Workers & Pages"**
3. 在列表中找到刚才创建的 Worker（如 `jeriblog-oauth-proxy`）
4. 点击 Worker 名称进入管理页面

#### 步骤 2：打开代码编辑器

在 Worker 管理页面中：
1. 找到右上角的 **"Edit code"** 按钮（此时应该是可点击的）
2. 点击 **"Edit code"** 按钮
3. 浏览器会打开一个新的代码编辑器页面

#### 步骤 3：替换代码

在代码编辑器中：
1. **删除所有默认代码**（左侧编辑区的全部内容）
2. 打开本地项目文件：`server/pkg/auth/cfw_proxy.js`
3. 复制完整内容（从第 1 行到最后一行）
4. 粘贴到 Cloudflare 代码编辑器的左侧编辑区
5. 右侧预览区会实时显示代码效果

#### 步骤 4：部署代码

1. 检查代码是否完整粘贴（确认有 `export default` 和完整的路由表）
2. 点击右上角的蓝色 **"Deploy"** 按钮
3. 等待部署完成（通常几秒钟）
4. 部署成功后会显示 "Deployment successful" 提示

**注意**：如果 Deploy 按钮是灰色的，说明代码没有修改或有语法错误，请检查代码是否完整。

### 3.4 获取 Worker 地址

部署成功后，Worker 地址会显示在多个位置：

#### 方式 1：从编辑器页面获取
- 代码编辑器顶部会显示 Worker 地址
- 格式：`https://<worker-name>.<your-username>.workers.dev`

#### 方式 2：从管理页面获取
1. 返回 Worker 管理页面
2. 在 "Preview" 区域可以看到完整地址
3. 或者在 "Settings" → "Triggers" 中查看 "Routes"

#### 地址示例

```
https://jeriblog-oauth-proxy.whizhzl.workers.dev
```

**地址组成说明**：
- `jeriblog-oauth-proxy`：您设置的 Worker 名称
- `whizhzl`：您的 Cloudflare 用户名
- `workers.dev`：Cloudflare Workers 的默认域名

**保存此地址**，后续配置项目时需要使用。


## 四、配置项目使用代理

### 4.1 修改 goth.go

编辑文件：`server/pkg/auth/goth.go`

找到第 30 行的 `workerProxy` 常量，修改为您的 Worker 地址：

```go
// workerProxy Cloudflare Worker 的代理地址
const workerProxy = "https://jeriblog-oauth-proxy.whizhzl.workers.dev"
```

### 4.2 更新 cfw_proxy.js 注释（可选）

编辑文件：`server/pkg/auth/cfw_proxy.js`

更新第 5 行的部署地址注释：

```javascript
/**
 * Cloudflare Worker - OAuth 代理
 * 用于代理 GitHub、Google 和 Microsoft 的 OAuth API 请求
 *
 * 部署地址: https://jeriblog-oauth-proxy.whizhzl.workers.dev
 *
 * 路由映射:
 *   /github-api/*  -> https://api.github.com/*
 *   /github/*      -> https://github.com/*
 *   ...
 */
```

## 五、测试代理服务

### 5.1 测试 GitHub API 代理

在浏览器中访问：
```
https://jeriblog-oauth-proxy.whizhzl.workers.dev/github-api/
```

如果返回 GitHub API 的响应（JSON 格式），说明代理工作正常。

### 5.2 测试 Google API 代理

在浏览器中访问：
```
https://jeriblog-oauth-proxy.whizhzl.workers.dev/google-api/
```

如果返回 Google API 的响应，说明代理工作正常。

### 5.3 测试根路径

访问根路径：
```
https://jeriblog-oauth-proxy.whizhzl.workers.dev/
```

会返回：
```json
{"error":"Invalid path"}
```

这是正常的，因为根路径没有匹配到任何路由规则。

## 六、常见问题

### 6.1 Worker 部署失败

**原因**：代码语法错误或格式问题

**解决方案**：
1. 检查是否完整复制了 `cfw_proxy.js` 的内容
2. 确保没有多余的字符或格式错误
3. 查看编辑器底部的错误提示

### 6.2 代理请求失败

**原因**：Worker 地址配置错误

**解决方案**：
1. 确认 `goth.go` 中的 `workerProxy` 地址正确
2. 确认 Worker 已成功部署
3. 测试 Worker 地址是否可以访问

### 6.3 OAuth 登录失败

**原因**：OAuth 配置不完整

**解决方案**：
1. 确认已在后台管理中配置 OAuth 的 Client ID 和 Secret
2. 确认回调地址配置正确
3. 检查代理服务是否正常工作

### 6.4 免费账号限制

Cloudflare Worker 免费账号限制：
- 每天 100,000 次请求
- 每次请求最多 10ms CPU 时间
- 每次请求最多 128MB 内存

对于个人博客来说，免费额度完全够用。

## 七、替代方案

### 7.1 使用原作者的代理（临时方案）

如果不想自己部署，可以暂时使用原作者的代理：

```go
const workerProxy = "https://jeriblog-oauth-proxy.whizhzl.workers.dev"
```

**注意**：不保证长期可用，建议部署自己的代理服务。

### 7.2 禁用 OAuth 第三方登录

如果不需要第三方登录功能：

1. 登录管理后台
2. 进入 `系统设置` → `OAuth 配置`
3. 不配置或清空 GitHub/Google/Microsoft 的 Client ID 和 Secret

这样用户只能使用邮箱密码登录。

## 八、安全建议

### 8.1 定期检查 Worker 日志

在 Cloudflare 控制台查看 Worker 的访问日志，监控异常请求。

### 8.2 限制请求来源（可选）

如果需要限制只允许特定域名访问，可以在 `cfw_proxy.js` 中添加来源检查：

```javascript
async fetch(request) {
  const origin = request.headers.get('Origin');
  const allowedOrigins = ['https://your-domain.com'];
  
  if (origin && !allowedOrigins.includes(origin)) {
    return new Response('Forbidden', { status: 403 });
  }
  
  // 原有代码...
}
```

### 8.3 监控使用量

定期检查 Worker 的请求量，避免超出免费额度。

## 九、总结

通过 Cloudflare Worker 部署 OAuth 代理服务，可以有效解决国内网络访问 GitHub、Google、Microsoft OAuth API 的问题。整个过程简单快捷，免费账号即可满足个人博客的需求。

**关键步骤回顾**：
1. ✅ 注册 Cloudflare 账号
2. ✅ 创建 Worker 并部署代码
3. ✅ 获取 Worker 地址
4. ✅ 修改项目配置文件
5. ✅ 测试代理服务
6. ✅ 配置 OAuth 应用

完成以上步骤后，您的博客系统就可以正常使用第三方 OAuth 登录功能了！
