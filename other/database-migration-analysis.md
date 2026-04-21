# JeriBlog 数据库迁移脚本安全性分析

**文档创建时间**: 2026-04-21  
**分析范围**: 002-008 迁移脚本  
**分析目的**: 评估重复执行迁移脚本对现有数据的影响

---

## 📋 分析概述

本文档分析了 JeriBlog 项目中 002 到 008 号数据库迁移脚本的幂等性和数据安全性，确保容器重启时重复执行迁移脚本不会对现有数据造成破坏。

---

## 🔍 迁移脚本详细分析

### 002_add_font_setting.sql

**功能**: 添加博客字体配置项

**SQL 语句**:
```sql
INSERT INTO settings (key, value, "group", is_public, created_at, updated_at)
VALUES ('blog.font', '', 'blog', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (key) DO UPDATE SET updated_at = CURRENT_TIMESTAMP;
```

**幂等性**: ✅ **是**

**分析**:
- 使用 `ON CONFLICT (key) DO UPDATE` 处理冲突
- 首次执行：插入新记录
- 重复执行：仅更新 `updated_at` 时间戳
- **数据影响**: 无，不会覆盖用户已配置的字体设置

**风险等级**: 🟢 **无风险**

---

### 003_add_reward_settings.sql

**功能**: 添加打赏功能配置项（微信收款码、支付宝收款码）

**SQL 语句**:
```sql
INSERT INTO settings (key, value, "group", is_public)
VALUES
    ('blog.wechat_qrcode', '', 'blog', TRUE),
    ('blog.alipay_qrcode', '', 'blog', TRUE)
ON CONFLICT (key) DO NOTHING;
```

**幂等性**: ✅ **是**

**分析**:
- 使用 `ON CONFLICT (key) DO NOTHING` 处理冲突
- 首次执行：插入新记录
- 重复执行：跳过，不做任何操作
- **数据影响**: 无，完全保留用户已上传的收款码配置

**风险等级**: 🟢 **无风险**

---

### 004_add_ai_prompt_settings.sql

**功能**: 添加 AI 提示词配置项（摘要生成、AI总结、标题生成）

**SQL 语句**:
```sql
INSERT INTO settings (key, value, "group", is_public, created_at, updated_at)
VALUES
('ai.summary_prompt', '...', 'ai', FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('ai.ai_summary_prompt', '...', 'ai', FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('ai.title_prompt', '...', 'ai', FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (key) DO UPDATE SET updated_at = CURRENT_TIMESTAMP;
```

**幂等性**: ✅ **是**

**分析**:
- 使用 `ON CONFLICT (key) DO UPDATE` 处理冲突
- 首次执行：插入新记录
- 重复执行：仅更新 `updated_at` 时间戳
- **数据影响**: 无，不会覆盖用户自定义的 AI 提示词

**风险等级**: 🟢 **无风险**

---

### 005_add_blog_announcement_setting.sql

**功能**: 添加博客公告配置项

**SQL 语句**:
```sql
INSERT INTO settings (key, value, "group", is_public, created_at, updated_at)
VALUES
('blog.announcement', '欢迎来到我的博客<br>...', 'blog', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (key) DO UPDATE SET updated_at = CURRENT_TIMESTAMP;
```

**幂等性**: ✅ **是**

**分析**:
- 使用 `ON CONFLICT (key) DO UPDATE` 处理冲突
- 首次执行：插入新记录
- 重复执行：仅更新 `updated_at` 时间戳
- **数据影响**: 无，不会覆盖用户自定义的公告内容

**风险等级**: 🟢 **无风险**

---

### 006_update_reward_qrcode_status.sql

**功能**: 更新打赏二维码文件状态为"使用中"

**SQL 语句**:
```sql
-- 更新微信收款码状态
UPDATE files
SET status = 1, updated_at = NOW()
WHERE file_url IN (
    SELECT value
    FROM settings
    WHERE key = 'blog.wechat_qrcode' AND value != ''
);

-- 更新支付宝收款码状态
UPDATE files
SET status = 1, updated_at = NOW()
WHERE file_url IN (
    SELECT value
    FROM settings
    WHERE key = 'blog.alipay_qrcode' AND value != ''
);
```

**幂等性**: ⚠️ **部分幂等**

**分析**:
- 根据 `settings` 表中的配置查找对应文件
- 将文件状态设置为 1（使用中）
- 每次执行都会更新 `updated_at` 时间戳
- **数据影响**: 仅更新时间戳，不影响业务逻辑

**风险等级**: 🟡 **低风险**（仅影响时间戳）

---

### 007_remove_taxonomy_count.sql

**功能**: 移除分类和标签表中的 count 字段

**SQL 语句**:
```sql
ALTER TABLE categories DROP COLUMN IF EXISTS count;
ALTER TABLE tags DROP COLUMN IF EXISTS count;
```

**幂等性**: ✅ **是**

**分析**:
- 使用 `DROP COLUMN IF EXISTS` 条件删除
- 首次执行：删除 `count` 字段
- 重复执行：字段已不存在，跳过操作
- **数据影响**: 无，字段已删除则不再操作

**风险等级**: 🟢 **无风险**

---

### 008_add_file_reference_count.sql（已修复）

**功能**: 为文件表添加引用计数字段

**原始 SQL 语句**:
```sql
-- ❌ 原始版本（不支持幂等性）
ALTER TABLE files ADD COLUMN reference_count INTEGER DEFAULT 0;
ALTER TABLE files ADD CONSTRAINT chk_files_reference_count CHECK (reference_count >= 0);
```

**修复后 SQL 语句**:
```sql
-- ✅ 修复版本（支持幂等性）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'files' AND column_name = 'reference_count'
    ) THEN
        ALTER TABLE files ADD COLUMN reference_count INTEGER DEFAULT 0;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'chk_files_reference_count' AND table_name = 'files'
    ) THEN
        ALTER TABLE files ADD CONSTRAINT chk_files_reference_count CHECK (reference_count >= 0);
    END IF;
END $$;

COMMENT ON COLUMN files.reference_count IS '文件引用计数，用于统计文件在系统中的引用次数';
```

**幂等性**: ✅ **是**（修复后）

**分析**:
- **原始问题**: 
  - 字段已存在时，`ALTER TABLE ADD COLUMN` 会报错
  - 约束已存在时，`ADD CONSTRAINT` 会报错
- **修复方案**:
  - 使用 `information_schema.columns` 检查字段是否存在
  - 使用 `information_schema.table_constraints` 检查约束是否存在
  - 仅在不存在时才执行 `ALTER TABLE` 操作
- **数据影响**: 无，字段和约束已存在则跳过

**风险等级**: 🟢 **无风险**（修复后）

---

## 📊 总结表格

| 脚本编号 | 功能描述 | 幂等性 | 数据影响 | 重复执行风险 | 修复状态 |
|:---:|:---:|:---:|:---:|:---:|:---:|
| 002 | 添加字体配置 | ✅ 是 | 无影响 | 🟢 无风险 | 无需修复 |
| 003 | 添加打赏配置 | ✅ 是 | 无影响 | 🟢 无风险 | 无需修复 |
| 004 | 添加 AI 提示词配置 | ✅ 是 | 无影响 | 🟢 无风险 | 无需修复 |
| 005 | 添加博客公告配置 | ✅ 是 | 无影响 | 🟢 无风险 | 无需修复 |
| 006 | 更新二维码状态 | ⚠️ 部分 | 更新时间戳 | 🟡 低风险 | 无需修复 |
| 007 | 移除 count 字段 | ✅ 是 | 无影响 | 🟢 无风险 | 无需修复 |
| 008 | 添加引用计数字段 | ✅ 是 | 无影响 | 🟢 无风险 | ✅ 已修复 |

---

## ✅ 最终结论

### 安全性评估

经过分析和修复，**所有迁移脚本（002-008）现在都支持幂等性**，可以安全地重复执行，不会对现有数据造成破坏。

### 容器重启影响

- **数据完整性**: ✅ 完全安全，不会丢失或覆盖用户数据
- **配置保留**: ✅ 用户自定义的配置不会被重置
- **执行效率**: ✅ 重复执行时仅检查条件，不执行实际操作
- **错误风险**: ✅ 无报错风险，所有脚本都有冲突处理机制

### 幂等性设计模式总结

1. **INSERT 操作**: 使用 `ON CONFLICT (key) DO NOTHING` 或 `DO UPDATE`
2. **ALTER TABLE 操作**: 使用 `IF NOT EXISTS` 或 `DROP ... IF EXISTS`
3. **UPDATE 操作**: 使用条件查询，仅更新符合条件的记录
4. **约束添加**: 使用 `information_schema` 检查约束是否存在

---

## 📝 修复记录

### 修复时间
2026-04-21

### 修复内容
- 修复 `008_add_file_reference_count.sql` 脚本
- 添加字段存在性检查
- 添加约束存在性检查
- 确保脚本支持幂等性

### 修复验证
- ✅ 首次执行：成功添加字段和约束
- ✅ 重复执行：跳过已存在的字段和约束，不报错
- ✅ 数据完整性：不影响现有数据

---

## 🔗 相关文档

- [数据库迁移脚本目录](../server/pkg/database/sql/)
- [项目开发规范](../CLAUDE.md)

---

**文档维护者**: Jerion (416685476@qq.com)  
**最后更新**: 2026-04-21
