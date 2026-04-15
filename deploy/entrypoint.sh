#!/bin/sh
#********************************************************************
# 项目名称：JeriBlog
# 文件名称：entrypoint.sh
# 创建时间：2026-03-31 11:50:09
#
# 系统用户：Administrator
# 作　　者：Jerion
# 联系邮箱：416685476@qq.com
# 功能描述：容器启动脚本，创建持久化目录、注入管理后台运行时配置，然后启动 supervisord
#********************************************************************

set -e

# 创建持久化数据目录（挂载卷时目录可能不存在）
mkdir -p /app/data/uploads /app/data/logs

# 生成管理后台运行时配置文件
# API_URL 在 docker run / compose 中通过环境变量注入，默认指向本机后端
cat > /usr/share/nginx/html/admin/config.js << EOF
window.__APP_CONFIG__ = {
  apiUrl: "${API_URL:-http://localhost:8080/api/v1}"
};
EOF

echo "[entrypoint] admin API_URL = ${API_URL:-http://localhost:8080/api/v1}"

# 以 supervisord 替换当前进程（保证信号正确传递）
exec /usr/bin/supervisord -n -c /etc/supervisor/conf.d/supervisord.conf
