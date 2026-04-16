-- 添加打赏配置项
-- 此脚本用于为已存在的数据库添加打赏功能所需的配置项

INSERT INTO settings (key, value, "group", is_public)
VALUES
    ('blog.wechat_qrcode', '', 'blog', TRUE),
    ('blog.alipay_qrcode', '', 'blog', TRUE)
ON CONFLICT (key) DO NOTHING;
