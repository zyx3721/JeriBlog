# 上传配置独立化迁移说明

## 概述

本次迁移将上传配置从共用字段模式改为每个存储类型独立配置的模式，解决了切换存储类型时配置混乱的问题。

## 问题描述

**迁移前的问题：**

所有存储类型（local/s3/cos/oss/kodo/r2/minio）共用同一组配置字段：
- `upload.access_key`
- `upload.secret_key`
- `upload.region`
- `upload.bucket`
- `upload.endpoint`
- `upload.domain`
- `upload.use_ssl`

这导致：
1. 切换存储类型时，之前的配置会被覆盖
2. 无法同时保存多个存储类型的配置
3. 切换回之前的存储类型需要重新配置

## 解决方案

**迁移后的结构：**

每个存储类型拥有独立的配置字段：

### Local 存储
- `upload.local.enabled`

### S3 存储
- `upload.s3.access_key`
- `upload.s3.secret_key`
- `upload.s3.region`
- `upload.s3.bucket`
- `upload.s3.endpoint`
- `upload.s3.domain`

### OSS 存储
- `upload.oss.access_key`
- `upload.oss.secret_key`
- `upload.oss.region`
- `upload.oss.bucket`
- `upload.oss.domain`

### COS 存储
- `upload.cos.secret_id`
- `upload.cos.secret_key`
- `upload.cos.region`
- `upload.cos.bucket`
- `upload.cos.domain`

### Kodo 存储
- `upload.kodo.access_key`
- `upload.kodo.secret_key`
- `upload.kodo.region`
- `upload.kodo.bucket`
- `upload.kodo.domain`

### R2 存储
- `upload.r2.access_key`
- `upload.r2.secret_key`
- `upload.r2.bucket`
- `upload.r2.endpoint`
- `upload.r2.domain`
- `upload.r2.use_ssl`

### MinIO 存储
- `upload.minio.access_key`
- `upload.minio.secret_key`
- `upload.minio.region`
- `upload.minio.bucket`
- `upload.minio.endpoint`
- `upload.minio.domain`
- `upload.minio.use_ssl`

## 迁移文件

### 1. 数据库迁移
- **文件**: `011_separate_upload_storage_configs.sql`
- **功能**: 
  - 读取当前存储类型和配置
  - 为每个存储类型创建独立配置字段
  - 将当前配置迁移到对应存储类型的字段中
  - 删除旧的共用配置字段

### 2. 初始化脚本更新
- **文件**: `001_init_database.sql`
- **功能**: 为新部署环境提供完整的独立配置字段

### 3. 测试脚本
- **文件**: `test_011_migration.sql`
- **功能**: 验证迁移是否成功

## 代码变更

### 后端变更

#### 1. 配置结构 (`server/config/config.go`)
```go
// 旧结构
type UploadConfig struct {
    StorageType string
    AccessKey   string
    SecretKey   string
    Region      string
    Bucket      string
    Endpoint    string
    Domain      string
    UseSSL      bool
}

// 新结构
type UploadConfig struct {
    StorageType string
    MaxFileSize int64
    PathPattern string
    Local       LocalStorageConfig
    S3          S3StorageConfig
    OSS         OSSStorageConfig
    COS         COSStorageConfig
    Kodo        KodoStorageConfig
    R2          R2StorageConfig
    MinIO       MinIOStorageConfig
}
```

#### 2. 配置加载 (`server/internal/service/setting.go`)
- 更新配置键常量
- 重写 `ApplyDatabaseConfig` 方法，加载所有存储类型的配置

#### 3. 存储实现 (`server/pkg/upload/storage/s3_unified.go`)
- 更新 `ensureClient` 方法，根据存储类型从对应配置结构获取参数
- 更新 `NewS3UnifiedStorage` 方法，根据存储类型获取 bucket

### 前端变更

#### 1. 配置表单 (`admin/src/views/file/components/UploadConfigDialog.vue`)
- 更新表单数据结构，为每个存储类型创建独立配置对象
- 添加 `currentConfig` 计算属性，根据当前存储类型返回对应配置
- 更新 `loadConfigs` 方法，加载所有存储类型的配置
- 更新 `handleSubmit` 方法，保存所有存储类型的配置

## 迁移步骤

### 对于已部署环境

1. **备份数据库**
   ```bash
   docker exec pg-prod pg_dump -U postgres jeriblog > backup_before_011.sql
   ```

2. **执行迁移脚本**
   ```bash
   docker exec -i pg-prod psql -U postgres -d jeriblog < server/pkg/database/sql/011_separate_upload_storage_configs.sql
   ```

3. **验证迁移**
   ```bash
   docker exec -i pg-prod psql -U postgres -d jeriblog < server/pkg/database/sql/test_011_migration.sql
   ```

4. **重启后端服务**
   ```bash
   docker-compose restart app
   ```

### 对于新部署环境

新部署环境会自动使用更新后的 `001_init_database.sql`，无需额外操作。

## 验证方法

### 1. 数据库验证
```sql
-- 检查旧字段是否已删除
SELECT key FROM settings WHERE key IN (
    'upload.access_key',
    'upload.secret_key',
    'upload.region',
    'upload.bucket',
    'upload.endpoint',
    'upload.domain',
    'upload.use_ssl'
);
-- 应返回 0 行

-- 检查新字段是否已创建
SELECT key FROM settings WHERE key LIKE 'upload.s3.%' ORDER BY key;
-- 应返回 6 行（S3 的 6 个配置字段）
```

### 2. 功能验证
1. 登录管理后台
2. 进入「文件管理」→「上传配置」
3. 切换不同的存储类型，验证：
   - 每个存储类型的配置独立保存
   - 切换回之前的存储类型，配置仍然保留
   - 保存配置后，上传功能正常

## 注意事项

1. **迁移不可逆**：执行迁移后，旧的共用配置字段会被删除
2. **配置保留**：迁移会将当前存储类型的配置迁移到对应字段，其他存储类型的字段为空
3. **热重载**：保存配置后会自动热重载，无需重启服务
4. **兼容性**：前后端必须同时更新，否则会出现配置读取错误

## 回滚方案

如果迁移后出现问题，可以通过以下步骤回滚：

1. **恢复数据库备份**
   ```bash
   docker exec -i pg-prod psql -U postgres -d jeriblog < backup_before_011.sql
   ```

2. **回滚代码**
   ```bash
   git revert <commit-hash>
   ```

3. **重启服务**
   ```bash
   docker-compose restart app
   ```

## 常见问题

### Q1: 迁移后上传功能报错？
**A**: 检查当前存储类型的配置是否完整，特别是 `access_key`、`secret_key`、`bucket` 等必需字段。

### Q2: 切换存储类型后配置丢失？
**A**: 这是正常现象。迁移只会保留当前存储类型的配置，其他存储类型需要重新配置。

### Q3: 前端显示配置为空？
**A**: 确保前后端代码都已更新，清除浏览器缓存后重试。

## 相关文件

- `server/pkg/database/sql/011_separate_upload_storage_configs.sql` - 迁移脚本
- `server/pkg/database/sql/001_init_database.sql` - 初始化脚本
- `server/pkg/database/sql/test_011_migration.sql` - 测试脚本
- `server/config/config.go` - 配置结构定义
- `server/internal/service/setting.go` - 配置加载逻辑
- `server/pkg/upload/storage/s3_unified.go` - 存储实现
- `admin/src/views/file/components/UploadConfigDialog.vue` - 前端配置界面

## 版本信息

- **迁移版本**: 011
- **目标版本**: v1.0.2
- **创建日期**: 2026-04-24
- **作者**: Jerion
