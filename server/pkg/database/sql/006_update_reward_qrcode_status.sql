/*
项目名称：JeriBlog
文件名称：006_update_reward_qrcode_status.sql
创建时间：2026-04-16 14:50:31

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：更新打赏二维码文件状态为"使用中"
*/

-- 更新微信收款码状态
UPDATE files
SET status = 1, updated_at = NOW()
WHERE file_url IN (
    SELECT value
    FROM settings
    WHERE `key` = 'blog.wechat_qrcode' AND value != ''
);

-- 更新支付宝收款码状态
UPDATE files
SET status = 1, updated_at = NOW()
WHERE file_url IN (
    SELECT value
    FROM settings
    WHERE `key` = 'blog.alipay_qrcode' AND value != ''
);
