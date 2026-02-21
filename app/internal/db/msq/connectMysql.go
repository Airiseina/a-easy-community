package msq

import (
	"commmunity/app/config"
	"commmunity/app/internal/model"
	"commmunity/app/zlog"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Gorm struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *Gorm {
	return &Gorm{
		db: db,
	}
}

func ConnectMysql() *gorm.DB {
	var sql config.MysqlConfig
	config.GetConfig()
	sql.Host = viper.GetString("mysql.host")
	sql.Port = viper.GetString("mysql.port")
	sql.User = viper.GetString("mysql.user")
	sql.Password = viper.GetString("mysql.password")
	sql.Name = viper.GetString("mysql.name")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", sql.User, sql.Password, sql.Host, sql.Port, sql.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zlog.Fatal("数据库连接失败", zap.Error(err))
	}
	err = db.AutoMigrate(&model.User{}, &model.UserProfile{}, &model.Post{}, &model.Comment{})
	if err != nil {
		zlog.Fatal("自动迁移失败", zap.Error(err))
	}
	zlog.Info("自动迁移成功")

	if !db.Migrator().HasIndex(&model.Post{}, "idx_fulltext_search") {
		err = db.Exec("ALTER TABLE posts ADD FULLTEXT INDEX idx_fulltext_search (title, content) WITH PARSER ngram").Error
		if err != nil {
			zlog.Error("创建全文索引失败", zap.Error(err))
		} else {
			zlog.Info("创建全文索引成功")
		}
	}

	return db
}
