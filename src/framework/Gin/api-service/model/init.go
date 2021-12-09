package model

import (
	"api-service/pkg/log"
	"fmt"
	"github.com/spf13/viper"

	// MySQL driver.
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Database struct {
	Self   *gorm.DB
	Docker *gorm.DB
}

var DB *Database

func openDB(username, password, addr, name string) *gorm.DB {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		true,
		//"Asia/Shanghai"),
		"Local")

	db, err := gorm.Open("mysql", config)
	if err != nil {
		log.Fatalf(err, "Database connection failed. Database name: %s", name)
	}
	// set for db connection
	setupDB(db)

	return db
}

func setupDB(db *gorm.DB) {
	db.LogMode(viper.GetBool("gormlog"))
	db.DB().SetMaxOpenConns(20000) // 设置最大的连接数，避免并发太高导致 too many connections。
	db.DB().SetMaxIdleConns(0)     // 设置闲置的连接数，当开启的一个连接使用完成后可以放在池里等候下一次使用。
}

// InitSelfDB used for cli
func InitSelfDB() *gorm.DB {
	return openDB(viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.addr"),
		viper.GetString("db.name"))
}

func InitDockerDB() *gorm.DB {
	return openDB(viper.GetString("docker_db.username"),
		viper.GetString("docker_db.password"),
		viper.GetString("docker_db.addr"),
		viper.GetString("docker_db.name"))
}

func (db *Database) Init() {
	DB = &Database{
		Self:   InitSelfDB(),
		Docker: InitDockerDB(),
	}
}

func (db *Database) Close() {
	_ = DB.Self.Close()
	_ = DB.Docker.Close()
}
