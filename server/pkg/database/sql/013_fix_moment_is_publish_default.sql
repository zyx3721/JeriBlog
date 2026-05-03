/*
项目名称：JeriBlog
文件名称：013_fix_moment_is_publish_default.sql
创建时间：2026-04-28 01:55:00

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：修复动态表 is_publish 字段默认值问题，移除默认值,允许前端明确控制发布状态
*/


-- 移除 is_publish 字段的默认值
ALTER TABLE moments ALTER COLUMN is_publish DROP DEFAULT;

-- 注释:移除默认值后,应用层必须明确传入 is_publish 的值
-- 这样可以避免 Go 零值(false)被 GORM 忽略的问题
