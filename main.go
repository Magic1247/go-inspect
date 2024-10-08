package main

import (
	"fmt"
	"go-inspect/config"
	"go-inspect/models"
	"go-inspect/routes"
	"go-inspect/tasks"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置文件
	viper.SetConfigFile("./config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 初始化数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.dbname"))
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	// 初始化全局 DB 变量
	config.InitDB(db)

	// 自动迁移模型
	config.DB.AutoMigrate(&models.User{}, &models.InspectionPoint{}, &models.InspectionRoute{}, &models.InspectionPlan{}, &models.InspectionOrder{}, &models.InspectionPointCheck{})

	// 设置 Gin 模式
	gin.SetMode(viper.GetString("server.mode"))

	// 创建 Gin 引擎
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r)

	// 初始化并启动定时任务
	tasks.InitCronJobs()

	// 启动服务器
	port := viper.GetString("server.port")
	fmt.Printf("Server is running on port %s\n", port)
	r.Run(":" + port)

	// 在程序结束时停止定时任务（这里可能需要适当的信号处理）
	defer tasks.StopCronJobs()
}
