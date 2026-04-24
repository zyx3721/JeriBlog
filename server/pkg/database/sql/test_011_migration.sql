/*
项目名称：JeriBlog
文件名称：test_011_migration.sql
创建时间：2026-04-24 10:45:00

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：测试 011 迁移脚本 - 验证上传配置迁移是否正确
*/

-- =============================================
-- 测试 011 迁移脚本
-- =============================================

DO $$
BEGIN
    RAISE NOTICE '========================================';
    RAISE NOTICE '开始测试 011 迁移脚本';
    RAISE NOTICE '========================================';

    -- 1. 检查旧配置字段是否已删除
    IF EXISTS (SELECT 1 FROM settings WHERE key IN (
        'upload.access_key',
        'upload.secret_key',
        'upload.region',
        'upload.bucket',
        'upload.endpoint',
        'upload.domain',
        'upload.use_ssl'
    )) THEN
        RAISE NOTICE '❌ 失败：旧配置字段仍然存在';
    ELSE
        RAISE NOTICE '✅ 通过：旧配置字段已删除';
    END IF;

    -- 2. 检查新配置字段是否已创建
    IF EXISTS (SELECT 1 FROM settings WHERE key LIKE 'upload.s3.%') AND
       EXISTS (SELECT 1 FROM settings WHERE key LIKE 'upload.oss.%') AND
       EXISTS (SELECT 1 FROM settings WHERE key LIKE 'upload.cos.%') AND
       EXISTS (SELECT 1 FROM settings WHERE key LIKE 'upload.kodo.%') AND
       EXISTS (SELECT 1 FROM settings WHERE key LIKE 'upload.r2.%') AND
       EXISTS (SELECT 1 FROM settings WHERE key LIKE 'upload.minio.%') THEN
        RAISE NOTICE '✅ 通过：所有存储类型的配置字段已创建';
    ELSE
        RAISE NOTICE '❌ 失败：部分存储类型的配置字段缺失';
    END IF;

    -- 3. 检查 Local 配置
    IF EXISTS (SELECT 1 FROM settings WHERE key = 'upload.local.enabled') THEN
        RAISE NOTICE '✅ 通过：Local 配置已创建';
    ELSE
        RAISE NOTICE '❌ 失败：Local 配置缺失';
    END IF;

    -- 4. 统计各存储类型的配置字段数量
    RAISE NOTICE '----------------------------------------';
    RAISE NOTICE '各存储类型配置字段统计：';
    RAISE NOTICE 'S3: % 个字段', (SELECT COUNT(*) FROM settings WHERE key LIKE 'upload.s3.%');
    RAISE NOTICE 'OSS: % 个字段', (SELECT COUNT(*) FROM settings WHERE key LIKE 'upload.oss.%');
    RAISE NOTICE 'COS: % 个字段', (SELECT COUNT(*) FROM settings WHERE key LIKE 'upload.cos.%');
    RAISE NOTICE 'Kodo: % 个字段', (SELECT COUNT(*) FROM settings WHERE key LIKE 'upload.kodo.%');
    RAISE NOTICE 'R2: % 个字段', (SELECT COUNT(*) FROM settings WHERE key LIKE 'upload.r2.%');
    RAISE NOTICE 'MinIO: % 个字段', (SELECT COUNT(*) FROM settings WHERE key LIKE 'upload.minio.%');

    -- 5. 显示当前存储类型及其配置
    DECLARE
        current_storage TEXT;
    BEGIN
        SELECT value INTO current_storage FROM settings WHERE key = 'upload.storage_type';
        RAISE NOTICE '----------------------------------------';
        RAISE NOTICE '当前存储类型: %', COALESCE(current_storage, 'local');

        IF current_storage IS NOT NULL AND current_storage != 'local' THEN
            RAISE NOTICE '当前存储类型的配置：';
            FOR rec IN
                SELECT key,
                       CASE
                           WHEN key LIKE '%secret%' OR key LIKE '%key%' THEN '***已隐藏***'
                           ELSE value
                       END as display_value
                FROM settings
                WHERE key LIKE 'upload.' || current_storage || '.%'
                ORDER BY key
            LOOP
                RAISE NOTICE '  %: %', rec.key, rec.display_value;
            END LOOP;
        END IF;
    END;

    RAISE NOTICE '========================================';
    RAISE NOTICE '测试完成';
    RAISE NOTICE '========================================';
END $$;
