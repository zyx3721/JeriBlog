-- 添加字体配置项
INSERT INTO settings (key, value, "group", is_public, created_at, updated_at)
VALUES ('blog.font', '', 'blog', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (key) DO UPDATE SET updated_at = CURRENT_TIMESTAMP;
