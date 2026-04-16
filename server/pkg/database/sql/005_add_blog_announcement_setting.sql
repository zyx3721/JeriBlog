/*
项目名称：JeriBlog
文件名称：005_add_blog_announcement_setting.sql
创建时间：2026-04-16 15:30:00

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：添加博客公告配置项
*/

-- 添加博客公告配置项
INSERT INTO settings (key, value, "group", is_public, created_at, updated_at)
VALUES
('blog.announcement', '欢迎来到我的博客<br><strong>记录技术、生活与思考</strong><br><span style="color:#e03131">内容持续更新中</span><br><a href="/about">了解更多</a>', 'blog', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (key) DO UPDATE SET updated_at = CURRENT_TIMESTAMP;
