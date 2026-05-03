

## 📊 邮件和飞书通知类型分析

### 📧 **邮件通知**（共 6 种）

| 序号 |     通知类型     |      模板文件       |     触发场景     |   接收者   | 使用地址 |
| :--: | :--------------: | :-----------------: | :--------------: | :--------: | :------: |
|  1   | `password_reset` | password_reset.tmpl | 用户请求重置密码 |  用户本人  | BlogURL  |
|  2   | `comment_reply`  | comment_reply.tmpl  |    评论被回复    |  被回复者  | BlogURL  |
|  3   |  `comment_new`   |  comment_new.tmpl   |    收到新评论    | 所有管理员 | AdminURL |
|  4   |  `feedback_new`  |  feedback_new.tmpl  |   收到反馈投诉   | 所有管理员 | AdminURL |
|  5   |  `friend_apply`  |  friend_apply.tmpl  |   收到友链申请   | 所有管理员 | AdminURL |
|  6   |    `default`     |    default.tmpl     |   默认通用模板   |     -      |    -     |

**邮件发送逻辑**：
- **文件位置**：`server/pkg/notification/email.go`
- **模板管理**：`server/pkg/email/template.go` (第115-129行)
- **地址选择**：根据 `isAdminNotification()` 判断使用 AdminURL 还是 BlogURL

---

### 🔔 **飞书通知**（共 4 种）

| 序号 |     通知类型     |      卡片构建方法      |   触发场景   |  接收者  |      使用地址      |
| :--: | :--------------: | :--------------------: | :----------: | :------: | :----------------: |
|  1   |  `comment_new`   |   buildCommentCard()   |  收到新评论  | 飞书群聊 | BlogURL + AdminURL |
|  2   |  `friend_apply`  | buildFriendApplyCard() | 收到友链申请 | 飞书群聊 |      AdminURL      |
|  3   |  `feedback_new`  |  buildFeedbackCard()   | 收到反馈投诉 | 飞书群聊 |      AdminURL      |
|  4   | `rss_feed_daily` |   buildRssFeedCard()   | RSS订阅日报  | 飞书群聊 |      AdminURL      |

**飞书发送逻辑**：
- **文件位置**：`server/pkg/notification/feishu.go`
- **卡片构建**：`buildCard()` 方法根据类型调用对应的卡片构建方法 (第47-60行)
- **特殊说明**：RSS日报通知仅发送飞书，不发送邮件和站内信

---

### 📱 **站内通知**（共 6 种）

| 序号 |     通知类型      |              通知方法              |   触发场景   |   接收者   |    显示位置    |
| :--: | :---------------: | :--------------------------------: | :----------: | :--------: | :------------: |
|  1   |  `comment_reply`  |        NotifyCommentReply()        |  评论被回复  |  被回复者  |   前台+后台    |
|  2   |   `comment_new`   |      NotifyCommentToAdmins()       |  收到新评论  | 所有管理员 |      后台      |
|  3   |  `feedback_new`   |          NotifyFeedback()          | 收到反馈投诉 | 所有管理员 |      后台      |
|  4   |  `friend_apply`   |        NotifyFriendApply()         | 收到友链申请 | 所有管理员 |      后台      |
|  5   | `friend_abnormal` |       NotifyFriendAbnormal()       | 友链异常提醒 | 所有管理员 | 后台（仅站内） |
|  6   |  `system_alert`   | NotifyVersionUpdateToSuperAdmins() | 版本更新提醒 | 超级管理员 | 后台（仅站内） |

**站内通知逻辑**：
- **文件位置**：`server/internal/service/notification.go`
- **存储表**：`notifications` + `user_notifications`
- **特殊说明**：
  - `friend_abnormal` 和 `system_alert` 仅发送站内信，不发送邮件和飞书
  - `rss_feed_daily` 仅发送飞书，不发送站内信和邮件

---

## 🎯 总结对比

|  通知渠道  | 通知类型数量 |                     特点                     |
| :--------: | :----------: | :------------------------------------------: |
|  **邮件**  |     6 种     |      包含密码重置、评论、反馈、友链申请      |
|  **飞书**  |     4 种     |      包含评论、反馈、友链申请、RSS日报       |
| **站内信** |     6 种     | 包含评论、反馈、友链申请、友链异常、系统通知 |

**通知渠道组合**：
- ✅ **三渠道都发**：`comment_new`、`feedback_new`、`friend_apply`
- ✅ **邮件+站内信**：`comment_reply`、`password_reset`
- ✅ **仅飞书**：`rss_feed_daily`
- ✅ **仅站内信**：`friend_abnormal`、`system_alert`

这样设计的好处是：
- 重要通知（评论、反馈、友链）多渠道覆盖，确保管理员及时收到
- 个人通知（密码重置、评论回复）仅发邮件和站内信，避免打扰
- 日常统计（RSS日报）仅发飞书，减少邮件干扰
- 系统提醒（友链异常、版本更新）仅站内信，避免过度通知