/*
项目名称：JeriBlog
文件名称：010_migrate_upload_type_for_existing_files.sql
创建时间：2026-04-23 14:30:00

系统用户：jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：为已部署环境迁移 upload_type 字段，将现有文件的引用关系补充到 upload_type 中（支持同用途多次引用）
*/

-- =============================================
-- 数据迁移说明
-- =============================================
-- 背景：upload_type 字段从单一用途改为支持多用途（逗号分隔）
-- 目标：为 reference_count > 0 的文件补充实际使用场景到 upload_type
-- 特性：支持同用途多次引用（如一个文件被 3 篇文章引用，则记录 3 次"文章配图"）
-- 用途类型：
--   用户头像、文章封面、文章配图、动态配图、动态视频、评论贴图、
--   友情链接A、友情链接S、站长头像、站长形象、博客图标、博客背景、
--   博客截图、展览图片、微信收款码、支付宝收款码、菜单图标、反馈投诉
-- =============================================

BEGIN;

-- 清空所有文件的 upload_type（重新计算）
UPDATE files SET upload_type = '';

-- 1. 用户头像（users.avatar）- 按引用次数追加
UPDATE files f
SET upload_type = (
    SELECT string_agg('用户头像', ',' ORDER BY u.id)
    FROM users u
    WHERE u.avatar = f.file_url
    AND u.deleted_at IS NULL
)
WHERE EXISTS (
    SELECT 1 FROM users u
    WHERE u.avatar = f.file_url
    AND u.deleted_at IS NULL
);

-- 2. 文章封面（articles.cover）- 按引用次数追加
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN (
        SELECT string_agg('文章封面', ',' ORDER BY a.id)
        FROM articles a
        WHERE a.cover = f.file_url
    )
    ELSE upload_type || ',' || (
        SELECT string_agg('文章封面', ',' ORDER BY a.id)
        FROM articles a
        WHERE a.cover = f.file_url
    )
END
WHERE EXISTS (
    SELECT 1 FROM articles a
    WHERE a.cover = f.file_url
);

-- 3. 文章配图（articles.content 中的图片）- 按引用次数追加
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN (
        SELECT string_agg('文章配图', ',' ORDER BY a.id)
        FROM articles a
        WHERE a.content LIKE '%' || f.file_url || '%'
    )
    ELSE upload_type || ',' || (
        SELECT string_agg('文章配图', ',' ORDER BY a.id)
        FROM articles a
        WHERE a.content LIKE '%' || f.file_url || '%'
    )
END
WHERE EXISTS (
    SELECT 1 FROM articles a
    WHERE a.content LIKE '%' || f.file_url || '%'
);

-- 4. 动态配图（moments.content 中的图片，JSON 格式）- 按引用次数追加
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN (
        SELECT string_agg('动态配图', ',' ORDER BY m.id)
        FROM moments m
        WHERE m.content::text LIKE '%' || f.file_url || '%'
        AND m.content::text LIKE '%"images"%'
    )
    ELSE upload_type || ',' || (
        SELECT string_agg('动态配图', ',' ORDER BY m.id)
        FROM moments m
        WHERE m.content::text LIKE '%' || f.file_url || '%'
        AND m.content::text LIKE '%"images"%'
    )
END
WHERE EXISTS (
    SELECT 1 FROM moments m
    WHERE m.content::text LIKE '%' || f.file_url || '%'
    AND m.content::text LIKE '%"images"%'
);

-- 5. 动态视频（moments.content 中的视频，JSON 格式）- 按引用次数追加
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN (
        SELECT string_agg('动态视频', ',' ORDER BY m.id)
        FROM moments m
        WHERE m.content::text LIKE '%' || f.file_url || '%'
        AND m.content::text LIKE '%"video"%'
    )
    ELSE upload_type || ',' || (
        SELECT string_agg('动态视频', ',' ORDER BY m.id)
        FROM moments m
        WHERE m.content::text LIKE '%' || f.file_url || '%'
        AND m.content::text LIKE '%"video"%'
    )
END
WHERE EXISTS (
    SELECT 1 FROM moments m
    WHERE m.content::text LIKE '%' || f.file_url || '%'
    AND m.content::text LIKE '%"video"%'
);

-- 6. 评论贴图（comments.content 中的图片）- 按引用次数追加
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN (
        SELECT string_agg('评论贴图', ',' ORDER BY c.id)
        FROM comments c
        WHERE c.content LIKE '%' || f.file_url || '%'
        AND c.deleted_at IS NULL
    )
    ELSE upload_type || ',' || (
        SELECT string_agg('评论贴图', ',' ORDER BY c.id)
        FROM comments c
        WHERE c.content LIKE '%' || f.file_url || '%'
        AND c.deleted_at IS NULL
    )
END
WHERE EXISTS (
    SELECT 1 FROM comments c
    WHERE c.content LIKE '%' || f.file_url || '%'
    AND c.deleted_at IS NULL
);

-- 7. 友情链接A（friends.avatar）- 按引用次数追加
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN (
        SELECT string_agg('友情链接A', ',' ORDER BY fr.id)
        FROM friends fr
        WHERE fr.avatar = f.file_url
    )
    ELSE upload_type || ',' || (
        SELECT string_agg('友情链接A', ',' ORDER BY fr.id)
        FROM friends fr
        WHERE fr.avatar = f.file_url
    )
END
WHERE EXISTS (
    SELECT 1 FROM friends fr
    WHERE fr.avatar = f.file_url
);

-- 8. 友情链接S（friends.screenshot）- 按引用次数追加
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN (
        SELECT string_agg('友情链接S', ',' ORDER BY fr.id)
        FROM friends fr
        WHERE fr.screenshot = f.file_url
    )
    ELSE upload_type || ',' || (
        SELECT string_agg('友情链接S', ',' ORDER BY fr.id)
        FROM friends fr
        WHERE fr.screenshot = f.file_url
    )
END
WHERE EXISTS (
    SELECT 1 FROM friends fr
    WHERE fr.screenshot = f.file_url
);

-- 9. 站长头像（settings.key = 'basic.author_avatar'）
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN '站长头像'
    ELSE upload_type || ',站长头像'
END
WHERE EXISTS (
    SELECT 1 FROM settings s
    WHERE s.value = f.file_url
    AND s.key = 'basic.author_avatar'
);

-- 10. 站长形象（settings.key = 'basic.author_photo'）
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN '站长形象'
    ELSE upload_type || ',站长形象'
END
WHERE EXISTS (
    SELECT 1 FROM settings s
    WHERE s.value = f.file_url
    AND s.key = 'basic.author_photo'
);

-- 11. 博客图标（settings.key = 'blog.favicon'）
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN '博客图标'
    ELSE upload_type || ',博客图标'
END
WHERE EXISTS (
    SELECT 1 FROM settings s
    WHERE s.value = f.file_url
    AND s.key = 'blog.favicon'
);

-- 12. 博客背景（settings.key = 'blog.background_image'）
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN '博客背景'
    ELSE upload_type || ',博客背景'
END
WHERE EXISTS (
    SELECT 1 FROM settings s
    WHERE s.value = f.file_url
    AND s.key = 'blog.background_image'
);

-- 13. 博客截图（settings.key = 'blog.screenshot'）
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN '博客截图'
    ELSE upload_type || ',博客截图'
END
WHERE EXISTS (
    SELECT 1 FROM settings s
    WHERE s.value = f.file_url
    AND s.key = 'blog.screenshot'
);

-- 14. 展览图片（settings.key = 'blog.about_exhibition'）
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN '展览图片'
    ELSE upload_type || ',展览图片'
END
WHERE EXISTS (
    SELECT 1 FROM settings s
    WHERE s.value = f.file_url
    AND s.key = 'blog.about_exhibition'
);

-- 15. 微信收款码（settings.key = 'blog.wechat_qrcode'）
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN '微信收款码'
    ELSE upload_type || ',微信收款码'
END
WHERE EXISTS (
    SELECT 1 FROM settings s
    WHERE s.value = f.file_url
    AND s.key = 'blog.wechat_qrcode'
);

-- 16. 支付宝收款码（settings.key = 'blog.alipay_qrcode'）
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN '支付宝收款码'
    ELSE upload_type || ',支付宝收款码'
END
WHERE EXISTS (
    SELECT 1 FROM settings s
    WHERE s.value = f.file_url
    AND s.key = 'blog.alipay_qrcode'
);

-- 17. 菜单图标（menus.icon）- 按引用次数追加
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN (
        SELECT string_agg('菜单图标', ',' ORDER BY m.id)
        FROM menus m
        WHERE m.icon = f.file_url
    )
    ELSE upload_type || ',' || (
        SELECT string_agg('菜单图标', ',' ORDER BY m.id)
        FROM menus m
        WHERE m.icon = f.file_url
    )
END
WHERE EXISTS (
    SELECT 1 FROM menus m
    WHERE m.icon = f.file_url
);

-- 18. 反馈投诉（feedbacks 表中的文件）- 按引用次数追加
UPDATE files f
SET upload_type = CASE
    WHEN upload_type = '' OR upload_type IS NULL THEN (
        SELECT string_agg('反馈投诉', ',' ORDER BY fb.id)
        FROM feedbacks fb
        WHERE fb.form_content LIKE '%' || f.file_url || '%'
    )
    ELSE upload_type || ',' || (
        SELECT string_agg('反馈投诉', ',' ORDER BY fb.id)
        FROM feedbacks fb
        WHERE fb.form_content LIKE '%' || f.file_url || '%'
    )
END
WHERE EXISTS (
    SELECT 1 FROM feedbacks fb
    WHERE fb.form_content LIKE '%' || f.file_url || '%'
);

-- 19. 对 upload_type 进行排序（保证一致性）
UPDATE files
SET upload_type = (
    SELECT string_agg(type_item, ',' ORDER BY type_item)
    FROM unnest(string_to_array(upload_type, ',')) AS type_item
)
WHERE upload_type != '' AND upload_type IS NOT NULL AND upload_type LIKE '%,%';

COMMIT;

-- =============================================
-- 验证查询（执行后可运行以下查询验证结果）
-- =============================================

-- 查看引用计数 > 0 但 upload_type 为空的文件（应该为 0）
-- SELECT id, file_name, file_url, reference_count, upload_type
-- FROM files
-- WHERE reference_count > 0 AND (upload_type = '' OR upload_type IS NULL);

-- 查看 upload_type 包含多个用途的文件
-- SELECT id, file_name, file_url, reference_count, upload_type
-- FROM files
-- WHERE upload_type LIKE '%,%'
-- ORDER BY reference_count DESC
-- LIMIT 20;

-- 统计各用途的文件数量（包含重复计数）
-- SELECT
--     unnest(string_to_array(upload_type, ',')) AS usage_type,
--     COUNT(*) AS usage_count
-- FROM files
-- WHERE upload_type != '' AND upload_type IS NOT NULL
-- GROUP BY usage_type
-- ORDER BY usage_count DESC;

-- 查看同用途多次引用的文件示例
-- SELECT id, file_name, file_url, reference_count, upload_type,
--        array_length(string_to_array(upload_type, ','), 1) AS usage_count
-- FROM files
-- WHERE upload_type LIKE '%文章配图,文章配图%'
--    OR upload_type LIKE '%用户头像,用户头像%'
-- ORDER BY reference_count DESC
-- LIMIT 10;
