package main

import (
	"fmt"
	"log"

	"github.com/subosito/gotenv"

	"jeri_blog/api/middleware"
	"jeri_blog/api/router"
	"jeri_blog/config"
	"jeri_blog/internal/service"
	"jeri_blog/pkg/database"
	"jeri_blog/pkg/logger"
	"jeri_blog/pkg/utils"

	_ "jeri_blog/docs" // swagger docs
)

// @title           JeriBlog
// @version         v1
// @description     一个基于 Go 语言的现代化博客后端服务

// @contact.name   Jerion
// @contact.email  416685476@qq.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 在请求头中添加 Bearer Token，格式：Bearer {token}

func main() {
	// 加载 .env 文件
	_ = gotenv.Load()

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	db, err := database.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		_ = db.Close()
	}()
	defer logger.Close()
	defer middleware.ClosePanicLogFile()

	// 执行数据库迁移
	if err := database.RunMigrations(db.DB); err != nil {
		middleware.ClosePanicLogFile()
		logger.Close()
		_ = db.Close()
		log.Fatalf("Failed to run migrations: %v", err) //nolint:gocritic // 已手动关闭资源
	}

	// 从数据库加载运行时配置（邮箱、上传等）
	settingService := service.NewSettingService(db.DB)
	settingService.SetConfig(cfg)               // 设置全局配置对象，用于热重载
	_ = settingService.ApplyDatabaseConfig(cfg) // 应用数据库配置

	// 初始化 IP 地理位置
	_ = utils.InitIPSearcher("")
	defer utils.CloseIPSearcher()

	// 初始化路由
	r := router.InitRouter(db, cfg)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("Server is running at http://%s:%d", cfg.Server.Host, cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
