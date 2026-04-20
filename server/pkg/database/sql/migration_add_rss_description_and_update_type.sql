-- RSS文章表添加 description 和 update_type 字段，修改唯一索引
-- 执行时间：2026-04-XX
-- 数据库：PostgreSQL

-- 1. 添加 description 字段（存储文章摘要/内容）
ALTER TABLE rss_articles ADD COLUMN IF NOT EXISTS description TEXT;

-- 2. 添加 update_type 字段（记录变更类型：content/title/published_at）
ALTER TABLE rss_articles ADD COLUMN IF NOT EXISTS update_type VARCHAR(20) DEFAULT '';

-- 3. 删除旧的 link 唯一索引（如果存在）
DROP INDEX IF EXISTS rss_articles_link_key;
DROP INDEX IF EXISTS idx_rss_articles_link_unique;

-- 4. 创建新的唯一索引：friend_id + title 组合
CREATE UNIQUE INDEX IF NOT EXISTS idx_rss_articles_friend_title ON rss_articles(friend_id, title);

-- 5. 为 link 字段创建普通索引（用于查询优化）
CREATE INDEX IF NOT EXISTS idx_rss_articles_link ON rss_articles(link);

-- 6. 添加字段注释
COMMENT ON COLUMN rss_articles.description IS '文章摘要/内容';
COMMENT ON COLUMN rss_articles.update_type IS '更新类型：content-内容已更新, title-标题已更新, published_at-发布时间已更新';

-- 说明：
-- 1. description 字段用于存储 RSS feed 的文章摘要或内容，用于检测内容变更
-- 2. update_type 字段记录文章的变更类型，前端会根据此字段显示不同颜色的标签
-- 3. 唯一索引从 link 改为 friend_id+title 组合，避免同一作者修改文章链接时产生重复文章
-- 4. 当文章内容/标题/发布时间变化时，会自动标记为未读并记录变更类型
