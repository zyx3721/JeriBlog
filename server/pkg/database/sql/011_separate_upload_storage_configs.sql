/*
项目名称：JeriBlog
文件名称：011_separate_upload_storage_configs.sql
创建时间：2026-04-24 10:30:00

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：分离上传配置 - 为每个存储类型创建独立的配置字段
*/

-- =============================================
-- 上传配置独立化迁移
-- =============================================
-- 目的：将共用的上传配置字段拆分为每个存储类型独立的配置
-- 影响：settings 表中的 upload 组配置
-- 版本：v1.0.2
-- =============================================

DO $$
DECLARE
    current_storage_type TEXT;
    current_access_key TEXT;
    current_secret_key TEXT;
    current_region TEXT;
    current_bucket TEXT;
    current_endpoint TEXT;
    current_domain TEXT;
    current_use_ssl TEXT;
BEGIN
    RAISE NOTICE '开始迁移上传配置...';

    -- 1. 读取当前配置
    SELECT value INTO current_storage_type FROM settings WHERE key = 'upload.storage_type';
    SELECT value INTO current_access_key FROM settings WHERE key = 'upload.access_key';
    SELECT value INTO current_secret_key FROM settings WHERE key = 'upload.secret_key';
    SELECT value INTO current_region FROM settings WHERE key = 'upload.region';
    SELECT value INTO current_bucket FROM settings WHERE key = 'upload.bucket';
    SELECT value INTO current_endpoint FROM settings WHERE key = 'upload.endpoint';
    SELECT value INTO current_domain FROM settings WHERE key = 'upload.domain';
    SELECT value INTO current_use_ssl FROM settings WHERE key = 'upload.use_ssl';

    RAISE NOTICE '当前存储类型: %', COALESCE(current_storage_type, 'local');

    -- 2. 确保基础配置字段存在（如果不存在则创建默认值）
    INSERT INTO settings (key, value, "group", is_public)
    VALUES ('upload.max_file_size', '10', 'upload', FALSE)
    ON CONFLICT (key) DO NOTHING;

    INSERT INTO settings (key, value, "group", is_public)
    VALUES ('upload.path_pattern', '{timestamp}_{random}{ext}', 'upload', FALSE)
    ON CONFLICT (key) DO NOTHING;

    -- 3. 为每个存储类型创建独立配置字段

    -- S3 存储配置
    INSERT INTO settings (key, value, "group", is_public) VALUES
    ('upload.s3.access_key', CASE WHEN current_storage_type = 's3' THEN COALESCE(current_access_key, '') ELSE '' END, 'upload', FALSE),
    ('upload.s3.secret_key', CASE WHEN current_storage_type = 's3' THEN COALESCE(current_secret_key, '') ELSE '' END, 'upload', FALSE),
    ('upload.s3.region', CASE WHEN current_storage_type = 's3' THEN COALESCE(current_region, '') ELSE '' END, 'upload', FALSE),
    ('upload.s3.bucket', CASE WHEN current_storage_type = 's3' THEN COALESCE(current_bucket, '') ELSE '' END, 'upload', FALSE),
    ('upload.s3.endpoint', CASE WHEN current_storage_type = 's3' THEN COALESCE(current_endpoint, '') ELSE '' END, 'upload', FALSE),
    ('upload.s3.domain', CASE WHEN current_storage_type = 's3' THEN COALESCE(current_domain, '') ELSE '' END, 'upload', FALSE)
    ON CONFLICT (key) DO NOTHING;

    -- OSS 存储配置
    INSERT INTO settings (key, value, "group", is_public) VALUES
    ('upload.oss.access_key', CASE WHEN current_storage_type = 'oss' THEN COALESCE(current_access_key, '') ELSE '' END, 'upload', FALSE),
    ('upload.oss.secret_key', CASE WHEN current_storage_type = 'oss' THEN COALESCE(current_secret_key, '') ELSE '' END, 'upload', FALSE),
    ('upload.oss.region', CASE WHEN current_storage_type = 'oss' THEN COALESCE(current_region, '') ELSE '' END, 'upload', FALSE),
    ('upload.oss.bucket', CASE WHEN current_storage_type = 'oss' THEN COALESCE(current_bucket, '') ELSE '' END, 'upload', FALSE),
    ('upload.oss.domain', CASE WHEN current_storage_type = 'oss' THEN COALESCE(current_domain, '') ELSE '' END, 'upload', FALSE)
    ON CONFLICT (key) DO NOTHING;

    -- COS 存储配置
    INSERT INTO settings (key, value, "group", is_public) VALUES
    ('upload.cos.secret_id', CASE WHEN current_storage_type = 'cos' THEN COALESCE(current_access_key, '') ELSE '' END, 'upload', FALSE),
    ('upload.cos.secret_key', CASE WHEN current_storage_type = 'cos' THEN COALESCE(current_secret_key, '') ELSE '' END, 'upload', FALSE),
    ('upload.cos.region', CASE WHEN current_storage_type = 'cos' THEN COALESCE(current_region, '') ELSE '' END, 'upload', FALSE),
    ('upload.cos.bucket', CASE WHEN current_storage_type = 'cos' THEN COALESCE(current_bucket, '') ELSE '' END, 'upload', FALSE),
    ('upload.cos.domain', CASE WHEN current_storage_type = 'cos' THEN COALESCE(current_domain, '') ELSE '' END, 'upload', FALSE)
    ON CONFLICT (key) DO NOTHING;

    -- Kodo 存储配置
    INSERT INTO settings (key, value, "group", is_public) VALUES
    ('upload.kodo.access_key', CASE WHEN current_storage_type = 'kodo' THEN COALESCE(current_access_key, '') ELSE '' END, 'upload', FALSE),
    ('upload.kodo.secret_key', CASE WHEN current_storage_type = 'kodo' THEN COALESCE(current_secret_key, '') ELSE '' END, 'upload', FALSE),
    ('upload.kodo.region', CASE WHEN current_storage_type = 'kodo' THEN COALESCE(current_region, '') ELSE '' END, 'upload', FALSE),
    ('upload.kodo.bucket', CASE WHEN current_storage_type = 'kodo' THEN COALESCE(current_bucket, '') ELSE '' END, 'upload', FALSE),
    ('upload.kodo.domain', CASE WHEN current_storage_type = 'kodo' THEN COALESCE(current_domain, '') ELSE '' END, 'upload', FALSE)
    ON CONFLICT (key) DO NOTHING;

    -- R2 存储配置
    INSERT INTO settings (key, value, "group", is_public) VALUES
    ('upload.r2.access_key', CASE WHEN current_storage_type = 'r2' THEN COALESCE(current_access_key, '') ELSE '' END, 'upload', FALSE),
    ('upload.r2.secret_key', CASE WHEN current_storage_type = 'r2' THEN COALESCE(current_secret_key, '') ELSE '' END, 'upload', FALSE),
    ('upload.r2.bucket', CASE WHEN current_storage_type = 'r2' THEN COALESCE(current_bucket, '') ELSE '' END, 'upload', FALSE),
    ('upload.r2.endpoint', CASE WHEN current_storage_type = 'r2' THEN COALESCE(current_endpoint, '') ELSE '' END, 'upload', FALSE),
    ('upload.r2.domain', CASE WHEN current_storage_type = 'r2' THEN COALESCE(current_domain, '') ELSE '' END, 'upload', FALSE),
    ('upload.r2.use_ssl', CASE WHEN current_storage_type = 'r2' THEN COALESCE(current_use_ssl, 'true') ELSE 'true' END, 'upload', FALSE)
    ON CONFLICT (key) DO NOTHING;

    -- MinIO 存储配置
    INSERT INTO settings (key, value, "group", is_public) VALUES
    ('upload.minio.access_key', CASE WHEN current_storage_type = 'minio' THEN COALESCE(current_access_key, '') ELSE '' END, 'upload', FALSE),
    ('upload.minio.secret_key', CASE WHEN current_storage_type = 'minio' THEN COALESCE(current_secret_key, '') ELSE '' END, 'upload', FALSE),
    ('upload.minio.region', CASE WHEN current_storage_type = 'minio' THEN COALESCE(current_region, '') ELSE '' END, 'upload', FALSE),
    ('upload.minio.bucket', CASE WHEN current_storage_type = 'minio' THEN COALESCE(current_bucket, '') ELSE '' END, 'upload', FALSE),
    ('upload.minio.endpoint', CASE WHEN current_storage_type = 'minio' THEN COALESCE(current_endpoint, '') ELSE '' END, 'upload', FALSE),
    ('upload.minio.domain', CASE WHEN current_storage_type = 'minio' THEN COALESCE(current_domain, '') ELSE '' END, 'upload', FALSE),
    ('upload.minio.use_ssl', CASE WHEN current_storage_type = 'minio' THEN COALESCE(current_use_ssl, 'true') ELSE 'true' END, 'upload', FALSE)
    ON CONFLICT (key) DO NOTHING;

    -- 3. 删除旧的共用配置字段
    DELETE FROM settings WHERE key IN (
        'upload.access_key',
        'upload.secret_key',
        'upload.region',
        'upload.bucket',
        'upload.endpoint',
        'upload.domain',
        'upload.use_ssl'
    );

    RAISE NOTICE '迁移完成！已将现有配置迁移到 % 存储类型', COALESCE(current_storage_type, 'local');
    RAISE NOTICE '旧的共用配置字段已删除';

END $$;
