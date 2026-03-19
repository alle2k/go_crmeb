package mysql

import (
	"boilerplate/config"
	"log"
	"strconv"
	"time"

	driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	cfg := driver.Config{
		User:      config.AppConfig.Mysql.Username,
		Passwd:    config.AppConfig.Mysql.Password,
		Net:       "tcp",
		Addr:      config.AppConfig.Mysql.Host + ":" + strconv.Itoa(config.AppConfig.Mysql.Port),
		DBName:    config.AppConfig.Mysql.Database,
		TLSConfig: "false",
		Params: map[string]string{
			"charset": "utf8mb4",
		},
		ParseTime:            true,
		Loc:                  time.Local,
		AllowNativePasswords: true,
	}
	dsn := cfg.FormatDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}
