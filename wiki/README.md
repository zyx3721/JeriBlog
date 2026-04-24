# JeriBlog Wiki 上传指南

本文档说明如何将这些 Wiki 页面上传到 GitHub。

## 前置条件

1. 确保你有 JeriBlog 仓库的写入权限
2. 已安装 Git 和 GitHub CLI（gh）

## 上传步骤

### 方法一：使用 GitHub 网页界面（推荐新手）

1. **初始化 Wiki**
   - 访问 https://github.com/zyx3721/JeriBlog/wiki
   - 点击 "Create the first page" 创建第一个页面
   - 标题填写 `Home`
   - 内容复制 `Home.md` 的内容
   - 点击 "Save Page"

2. **批量上传其他页面**
   - 点击右侧 "New Page" 按钮
   - 按照以下顺序创建页面：

   **快速入门**
   - 页面名：`Quick-Start`
   - 内容：复制 `Quick-Start.md`

   **总览页面**
   - 页面名：`Basic-Configuration`，内容：`Basic-Configuration.md`
   - 页面名：`User-Guide`，内容：`User-Guide.md`
   - 页面名：`FAQ`，内容：`FAQ.md`

   **基础配置（7个页面）**
   - `Basic-Config` ← `Basic-Config.md`
   - `Blog-Config` ← `Blog-Config.md`
   - `Upload-Config` ← `Upload-Config.md`
   - `OAuth-Config` ← `OAuth-Config.md`
   - `AI-Config` ← `AI-Config.md`
   - `Notification-Config` ← `Notification-Config.md`
   - `Import-Export` ← `Import-Export.md`

   **使用指南（7个页面）**
   - `Article-Management` ← `Article-Management.md`
   - `Article-Editor` ← `Article-Editor.md`
   - `Moment-Management` ← `Moment-Management.md`
   - `Friend-Link-Management` ← `Friend-Link-Management.md`
   - `Menu-Management` ← `Menu-Management.md`
   - `User-Management` ← `User-Management.md`

### 方法二：使用 Git 命令行（推荐熟练用户）

1. **初始化 Wiki（首次需要在网页创建第一个页面）**
   ```bash
   # 访问 https://github.com/zyx3721/JeriBlog/wiki
   # 点击 "Create the first page" 创建 Home 页面
   ```

2. **克隆 Wiki 仓库**
   ```bash
   cd /tmp
   git clone https://github.com/zyx3721/JeriBlog.wiki.git
   cd JeriBlog.wiki
   ```

3. **复制所有 Wiki 文件**
   ```bash
   cp "D:/GitHub/我的项目/JeriBlog/wiki"/*.md .
   ```

4. **提交并推送**
   ```bash
   git add *.md
   git commit -m "docs: 初始化 JeriBlog Wiki 文档"
   git push origin master
   ```

### 方法三：使用自动化脚本

运行以下命令自动上传所有页面：

```bash
cd "D:/GitHub/我的项目/JeriBlog/wiki"
bash upload_wiki.sh
```

## Wiki 页面结构

```
JeriBlog Wiki
├── Home.md                      # 主页
├── Quick-Start.md               # 快速入门
├── Basic-Configuration.md       # 基础配置总览
│   ├── Basic-Config.md          # 基本配置
│   ├── Blog-Config.md           # 博客配置
│   ├── Upload-Config.md         # 上传配置
│   ├── OAuth-Config.md          # OAuth 配置
│   ├── AI-Config.md             # AI 配置
│   ├── Notification-Config.md   # 通知配置
│   └── Import-Export.md         # 导入导出
├── User-Guide.md                # 使用指南总览
│   ├── Article-Management.md    # 文章管理
│   ├── Article-Editor.md        # 文章编辑
│   ├── Moment-Management.md     # 动态管理
│   ├── Friend-Link-Management.md # 友链管理
│   ├── Menu-Management.md       # 菜单管理
│   └── User-Management.md       # 用户管理
└── FAQ.md                       # 常见问题
```

## 页面导航

所有页面都包含了相互链接，用户可以方便地在页面间跳转：

- 每个详细页面底部都有"返回主页"、"返回总览"等导航链接
- 主页提供了完整的目录结构
- 总览页面提供了分类导航

## 验证上传

上传完成后，访问以下链接验证：

- Wiki 主页：https://github.com/zyx3721/JeriBlog/wiki
- 快速入门：https://github.com/zyx3721/JeriBlog/wiki/Quick-Start
- 基础配置：https://github.com/zyx3721/JeriBlog/wiki/Basic-Configuration
- 使用指南：https://github.com/zyx3721/JeriBlog/wiki/User-Guide
- 常见问题：https://github.com/zyx3721/JeriBlog/wiki/FAQ

## 后续维护

### 更新页面

1. 修改本地 Markdown 文件
2. 在 GitHub Wiki 页面点击 "Edit" 更新内容
3. 或使用 Git 命令推送更新

### 添加新页面

1. 在本地创建新的 Markdown 文件
2. 在相关页面添加链接
3. 上传到 GitHub Wiki

### 删除页面

1. 在 GitHub Wiki 页面点击 "Delete Page"
2. 或使用 Git 命令删除文件并推送

## 注意事项

1. **页面命名**：使用英文，单词间用连字符分隔，首字母大写
2. **链接格式**：使用 `[显示文本](页面名)` 格式，不需要 `.md` 后缀
3. **图片上传**：Wiki 支持拖拽上传图片，会自动生成链接
4. **Markdown 语法**：GitHub Wiki 支持标准 Markdown 和 GitHub Flavored Markdown

## 常见问题

**Q: Wiki 仓库克隆失败？**  
A: 需要先在 GitHub 网页上创建第一个 Wiki 页面，才能克隆仓库。

**Q: 页面链接失效？**  
A: 检查页面名称是否正确，注意大小写和连字符。

**Q: 如何设置侧边栏？**  
A: 创建 `_Sidebar.md` 文件，内容为导航链接列表。

**Q: 如何设置页脚？**  
A: 创建 `_Footer.md` 文件，内容会显示在每个页面底部。

## 技术支持

如有问题，请：

1. 查看 [GitHub Wiki 官方文档](https://docs.github.com/en/communities/documenting-your-project-with-wikis)
2. 在 [JeriBlog Issues](https://github.com/zyx3721/JeriBlog/issues) 提问
3. 联系项目维护者

---

📝 文档生成时间：2024-01-XX  
📦 共 18 个页面  
✨ 祝使用愉快！
