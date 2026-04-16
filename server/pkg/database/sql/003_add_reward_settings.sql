/*
项目名称：JeriBlog
文件名称：003_add_reward_settings.sql
创建时间：2026-04-16 15:30:00

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：添加打赏功能配置项（微信收款码、支付宝收款码）
*/

-- 添加打赏配置项

INSERT INTO settings (key, value, "group", is_public)
VALUES
    ('blog.wechat_qrcode', '', 'blog', TRUE),
    ('blog.alipay_qrcode', '', 'blog', TRUE)
ON CONFLICT (key) DO NOTHING;
