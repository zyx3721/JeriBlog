/*
项目名称：JeriBlog
文件名称：009_add_target_id_to_comments.sql
创建时间：2026-04-22 21:55:16

系统用户：jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：为 comments 表添加 target_id 字段，解决文章 slug 变更后评论来源显示问题
*/

-- 添加 target_id 字段（幂等性：仅在字段不存在时添加）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'comments' AND column_name = 'target_id'
    ) THEN
        ALTER TABLE comments ADD COLUMN target_id INTEGER;
    END IF;
END $$;

-- 添加索引（幂等性：仅在索引不存在时创建）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE tablename = 'comments' AND indexname = 'idx_comments_target_id'
    ) THEN
        CREATE INDEX idx_comments_target_id ON comments(target_id);
    END IF;
END $$;

-- 为现有评论填充 target_id（通过 slug 查询文章 ID）
UPDATE comments c
SET target_id = a.id
FROM articles a
WHERE c.target_type = 'article'
  AND c.target_key = a.slug
  AND c.target_id IS NULL;

-- 添加注释（幂等性：COMMENT 命令本身是幂等的）
COMMENT ON COLUMN comments.target_id IS '文章ID（用于解决slug变更问题）';
