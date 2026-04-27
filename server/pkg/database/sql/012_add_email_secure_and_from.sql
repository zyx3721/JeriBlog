/*
项目名称：JeriBlog
文件名称：012_add_email_secure_and_from.sql
创建时间：2026-04-28 01:55:00

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：添加邮件加密方式和发件人邮箱配置项
*/

-- 添加邮件加密方式配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('notification.email_secure', 'ssl', 'notification', FALSE)
ON CONFLICT (key) DO NOTHING;

-- 添加发件人邮箱配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('notification.email_from', '', 'notification', FALSE)
ON CONFLICT (key) DO NOTHING;
