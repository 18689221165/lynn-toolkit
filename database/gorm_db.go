package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
	"youliao.cn/liaoma-toolkit/types"
)

// Conf 数据库相关配置
type Conf struct {
	Name                      string `yaml:"name"`
	Dsn                       string `yaml:"dsn"`
	MaxOpenConn               int    `yaml:"maxOpenConn"`
	MaxIdleConn               int    `yaml:"maxIdleConn"`
	IgnoreRecordNotFoundError bool   `yaml:"ignoreRecordNotFoundError"`
	LogLevel                  int    `yaml:"logLevel"`
}

// NewDBClient 创建GORM数据库连接池
func NewDBClient(conf Conf) *gorm.DB {
	newlog := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second * 3,
			Colorful:                  false,
			IgnoreRecordNotFoundError: conf.IgnoreRecordNotFoundError,
			LogLevel:                  logger.LogLevel(conf.LogLevel),
		},
	)

	config := &gorm.Config{
		Logger: newlog,
		// 使用单数表明
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	db, err := gorm.Open(mysql.Open(conf.Dsn), config)
	if err != nil {
		log.Fatalf("数据库连接初始化失败,DSN:%s,错误：%+v", conf.Dsn, err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(conf.MaxOpenConn)
	sqlDB.SetMaxIdleConns(conf.MaxIdleConn)
	sqlDB.SetConnMaxIdleTime(time.Hour)
	sqlDB.SetConnMaxLifetime(2 * time.Hour)

	if err = sqlDB.Ping(); err != nil {
		log.Panicf("数据库连接失败: %+v", err)
	}
	return db
}

func NewDBClientWithProfile(profile types.Profile, conf Conf) *gorm.DB {
	c := NewDBClient(conf)
	if profile == types.Profile_Dev {
		c = c.Debug()
	}
	return c
}

func NewDBClientSetWithProfile(profile types.Profile, cfgset []Conf) map[string]*gorm.DB {
	dbmap := map[string]*gorm.DB{}
	for _, conf := range cfgset {
		dbclient := NewDBClientWithProfile(profile, conf)
		dbmap[conf.Name] = dbclient
	}
	return dbmap
}
