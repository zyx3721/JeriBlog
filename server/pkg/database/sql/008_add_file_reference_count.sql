/*
项目名称：JeriBlog
文件名称：008_add_file_reference_count.sql
创建时间：2026-04-21 10:30:00

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：为文件表添加引用计数字段，实现精准的文件使用状态判断
*/

-- =============================================
-- 添加引用计数字段
-- =============================================

-- 1. 添加 reference_count 字段（如果不存在）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'files' AND column_name = 'reference_count'
    ) THEN
        ALTER TABLE files ADD COLUMN reference_count INTEGER DEFAULT 0;
    END IF;
END $$;

-- 2. 添加约束检查（如果不存在）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'chk_files_reference_count' AND table_name = 'files'
    ) THEN
        ALTER TABLE files ADD CONSTRAINT chk_files_reference_count CHECK (reference_count >= 0);
    END IF;
END $$;

-- 3. 为引用计数字段添加注释
COMMENT ON COLUMN files.reference_count IS '文件引用计数，用于统计文件在系统中的引用次数';

-- =============================================
-- 初始化现有数据的引用计数
-- =============================================

-- 根据当前状态初始化引用计数（手动执行）
-- status = 1 (使用中) 的文件设置为 1，status = 0 (未使用) 的文件保持为 0
-- UPDATE files SET reference_count = CASE WHEN status = 1 THEN 1 ELSE 0 END WHERE reference_count = 0;
