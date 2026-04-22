/*
项目名称：JeriBlog
文件名称：20260422_add_target_id_to_comments.sql
创建时间：2026-04-22 21:55:16

系统用户：jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：为 comments 表添加 target_id 字段，解决文章 slug 变更后评论来源显示问题
*/

-- 添加 target_id 字段
ALTER TABLE comments ADD COLUMN target_id INTEGER;

-- 添加索引
CREATE INDEX idx_comments_target_id ON comments(target_id);

-- 为现有评论填充 target_id（通过 slug 查询文章 ID）
UPDATE comments c
SET target_id = a.id
FROM articles a
WHERE c.target_type = 'article'
  AND c.target_key = a.slug
  AND c.target_id IS NULL;

-- 添加注释
COMMENT ON COLUMN comments.target_id IS '文章ID（用于解决slug变更问题）';
