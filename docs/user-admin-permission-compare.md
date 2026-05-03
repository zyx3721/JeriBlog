## 超级管理员(super_admin)与管理员(admin)权限对比

### 1. 核心权限差异

| 功能模块           | 超级管理员       | 管理员               | 说明                             |
| ------------------ | ---------------- | -------------------- | -------------------------------- |
| **系统设置修改**   | ✅ 完全权限       | ❌ 无权限             | 仅超管可修改基础设置、通知设置等 |
| **AI MCP密钥重置** | ✅ 可重置         | ❌ 无权限             | 仅超管可重置MCP Secret           |
| **RSS订阅管理**    | ✅ 完全权限       | ❌ 无权限             | 标记已读、刷新订阅仅超管可操作   |
| **访问日志管理**   | ✅ 完全权限       | ⚠️ 仅查看权限         | 管理员仅可查看，无法删除日志     |
| **用户管理**       | ✅ 管理所有用户   | ⚠️ 仅管理普通用户     | 管理员无法操作超管和其他管理员   |
| **角色分配**       | ✅ 可分配任何角色 | ⚠️ 仅可分配user/guest | 管理员无法创建或提升为管理员     |
| **删除用户**       | ✅ 可删除任何用户 | ⚠️ 仅可删除普通用户   | 管理员无法删除超管和其他管理员   |
| **密码重置**       | ❌ 禁止邮箱重置   | ❌ 禁止邮箱重置       | 两者均不可通过验证码重置密码     |

### 2. API路由访问限制

**仅超级管理员可访问的路由**:

```go
// 系统设置组
PATCH /settings/:group                    // 更新设置组
PUT   /settings/ai/mcp-secret/reset      // 重置MCP密钥

// RSS订阅管理
PUT   /rssfeed/:id/read                  // 标记单篇文章已读
PUT   /rssfeed/read-all                  // 标记所有文章已读
POST  /rssfeed/refresh                   // 刷新RSS订阅源
```

**管理员及以上可访问的路由**:
- 所有其他管理功能(评论管理、文章管理、友链管理、反馈管理等)

### 3. 访问日志管理权限细节

#### 超级管理员权限:
- ✅ 查看所有访问日志
- ✅ 删除单条访问日志
- ✅ 批量删除访问日志
- ✅ 可选择日志记录(显示选择框)

#### 管理员权限:
- ✅ 查看所有访问日志
- ❌ 无法删除访问日志(操作列不显示删除按钮)
- ❌ 无法批量删除(不显示批量删除按钮)
- ❌ 无法选择日志记录(不显示选择框)

**前端实现** (`admin/src/views/visit/VisitList.vue`):
```vue
<!-- 选择列 (仅超级管理员可见) -->
<el-table-column
  v-if="isSuperAdmin"
  type="selection"
  width="55"
  align="center"
/>

<!-- 操作列 (仅超级管理员显示删除按钮) -->
<el-table-column label="操作" width="90" align="center" fixed="right">
  <template #default="{ row }">
    <template v-if="isSuperAdmin">
      <el-button type="danger" link @click="handleDelete(row.id)">删除</el-button>
    </template>
  </template>
</el-table-column>

<!-- 批量删除按钮 (仅超级管理员可见) -->
<el-button
  v-if="isSuperAdmin && selectedIds.length > 0"
  type="danger"
  @click="handleBatchDelete"
>
  批量删除
</el-button>
```

### 4. 用户管理权限细节

#### 超级管理员权限:
- ✅ 查看所有用户(包括其他超管和管理员)
- ✅ 创建任何角色的用户
- ✅ 修改任何用户的角色
- ✅ 删除任何用户(除最后一个超管外)
- ✅ 禁用/启用任何用户

#### 管理员权限:
- ✅ 查看所有用户
- ⚠️ 仅可创建user/guest角色用户
- ⚠️ 仅可修改user/guest角色用户
- ⚠️ 仅可删除user/guest角色用户
- ⚠️ 无法操作超管和其他管理员账户

**代码实现** (`server/internal/service/user.go:811-869`):
```go
// 管理员权限检查
if currentUser.Role == model.RoleAdmin {
    // 管理员只能管理普通用户
    if targetUser.Role == model.RoleSuperAdmin || targetUser.Role == model.RoleAdmin {
        return errors.New("管理员无权操作其他管理员或超级管理员")
    }
}
```

**前端实现** (`admin/src/views/user/UserList.vue`):
```vue
<!-- 操作列 (根据权限显示按钮) -->
<el-table-column label="操作" width="180" align="center" fixed="right">
  <template #default="{ row }">
    <template v-if="!row.deleted_at && canOperateUser(row)">
      <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
      <el-button type="danger" link size="small" @click="handleDelete(row.id)">删除</el-button>
    </template>
  </template>
</el-table-column>

<script setup>
// 判断是否为管理角色
const isManagedRole = (role: string) => role === 'admin' || role === 'super_admin';

// 判断是否可以操作该用户
const canOperateUser = (user: UserType) =>
  currentRole.value === 'super_admin' || !isManagedRole(user.role);
</script>
```

### 5. 账户保护机制

#### 超级管理员保护:
- **最后一个超管保护**: 系统至少保留一个超级管理员
- **删除限制**: 无法删除最后一个超管账户
- **降级限制**: 无法将最后一个超管降级为其他角色

**代码实现** (`server/internal/service/user.go:729-790`):
```go
func (s *UserService) CountSuperAdmins() (int64, error) {
    return s.repo.CountByRole(model.RoleSuperAdmin)
}

// 删除前检查
if user.Role == model.RoleSuperAdmin {
    count, _ := s.CountSuperAdmins()
    if count <= 1 {
        return errors.New("无法删除最后一个超级管理员")
    }
}
```

### 6. 通知接收差异

| 通知类型         | 超级管理员 | 管理员   | 说明                   |
| ---------------- | ---------- | -------- | ---------------------- |
| **新评论通知**   | ✅ 接收     | ✅ 接收   | 所有管理员均接收       |
| **友链申请通知** | ✅ 接收     | ✅ 接收   | 所有管理员均接收       |
| **反馈投诉通知** | ✅ 接收     | ✅ 接收   | 所有管理员均接收       |
| **版本更新通知** | ✅ 接收     | ❌ 不接收 | 仅超管接收版本更新提醒 |

**代码实现** (`server/internal/repository/notification.go:136-153`):
```go
// 获取所有管理员(包括超管和管理员)
func (r *NotificationRepository) GetAllAdmins() ([]*model.User, error) {
    return r.db.Where("role IN ?", []string{model.RoleAdmin, model.RoleSuperAdmin}).Find(&users).Error
}

// 仅获取超级管理员
func (r *NotificationRepository) GetAllSuperAdmins() ([]*model.User, error) {
    return r.db.Where("role = ?", model.RoleSuperAdmin).Find(&users).Error
}
```

### 7. 中间件权限控制

**超级管理员专用中间件** (`server/api/middleware/rbac.go:43-86`):
```go
func IsSuperAdmin() gin.HandlerFunc {
    return func(c *gin.Context) {
        if user.Role != model.RoleSuperAdmin {
            c.JSON(http.StatusForbidden, gin.H{"error": "需要超级管理员权限"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

**管理员及以上中间件**:
```go
func IsAdminOrAbove() gin.HandlerFunc {
    return func(c *gin.Context) {
        if user.Role != model.RoleAdmin && user.Role != model.RoleSuperAdmin {
            c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 8. 特殊限制

#### 密码重置限制 (`server/internal/service/verification.go:83-85`):
```go
// 管理员和超级管理员均不可通过邮箱验证码重置密码
if user.Role == model.RoleAdmin || user.Role == model.RoleSuperAdmin {
    return errors.New("管理员账户不支持通过邮箱重置密码")
}
```

### 9. 权限绕过机制

**超级管理员权限绕过** (`server/api/middleware/rbac.go:43-86`):
- 如果用户角色为`super_admin`,所有权限检查自动通过
- 无需额外的权限验证逻辑

---

## 总结

**超级管理员**是系统的最高权限角色,拥有:
- ✅ 系统设置完全控制权
- ✅ RSS订阅管理权限
- ✅ 访问日志完全管理权限(查看、删除、批量删除)
- ✅ 用户管理完全权限
- ✅ 版本更新通知接收
- ✅ AI功能配置权限

**管理员**是受限的管理角色,仅拥有:
- ✅ 内容管理权限(评论、文章、友链、反馈)
- ⚠️ 受限的用户管理权限(仅普通用户)
- ⚠️ 访问日志查看权限(无删除权限)
- ❌ 无系统设置修改权限
- ❌ 无RSS订阅管理权限
- ❌ 无版本更新通知

**核心设计理念**:
- 超级管理员负责系统级配置和高级功能
- 管理员负责日常内容审核和用户管理
- 系统至少保留一个超级管理员账户
- 两者均不可通过邮箱验证码重置密码(安全考虑)