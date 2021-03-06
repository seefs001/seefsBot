package orm

import (
	"log"
	"time"

	"gorm.io/gorm/schema"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// db database instance
var dbInstance *gorm.DB

func DB() *gorm.DB {
	return dbInstance
}

// Init initialize database
func Init(dialector gorm.Dialector) error {
	gormLogger := logger.New(
		log.Default(), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 彩色打印
		},
	)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
		// disable foreign key constraint
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: &schema.NamingStrategy{
			TablePrefix: "tb_",
			// table name singular
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}
	dbInstance = db
	s, err := db.DB()
	if err != nil {
		return err
	}
	s.SetMaxIdleConns(50)
	s.SetMaxOpenConns(100)
	s.SetConnMaxLifetime(30 * time.Second)
	return nil
}

func Close() error {
	db, err := dbInstance.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
