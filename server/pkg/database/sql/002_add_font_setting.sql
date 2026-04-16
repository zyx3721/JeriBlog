/*
项目名称：JeriBlog
文件名称：002_add_font_setting.sql
创建时间：2026-04-16 15:30:00

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：添加博客字体配置项
*/

-- 添加字体配置项
INSERT INTO settings (key, value, "group", is_public, created_at, updated_at)
VALUES ('blog.font', '', 'blog', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (key) DO UPDATE SET updated_at = CURRENT_TIMESTAMP;
