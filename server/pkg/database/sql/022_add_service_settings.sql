-- 添加 Meting-API、头像服务、IP 查询、封面制作、OAuth 代理配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('blog.meting_api', 'https://api.injahow.cn/meting/', 'blog', TRUE),
('blog.cravatar_url', 'https://cravatar.cn/avatar/%s?s=200&d=robohash', 'blog', TRUE),
('blog.ip_api_url', 'http://ip-api.com/json/%s?lang=zh-CN', 'blog', TRUE),
('blog.cover_maker_api', 'https://pixhub.flec.top', 'blog', TRUE),
('oauth.worker_proxy', 'https://jeriblog-oauth-proxy.whizhzl.workers.dev', 'oauth', FALSE)
ON CONFLICT (key) DO NOTHING;
