-- 添加公众号二维码配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('blog.wechat_offaccounts', '', 'blog', TRUE),
('blog.wechat_offname', '', 'blog', TRUE)
ON CONFLICT (key) DO NOTHING;
