package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"cli_sederhana2/model" // Sesuaikan dengan jalur paket model Anda
)

type MysqlDB struct {
	DB *gorm.DB
}

var Mysql MysqlDB

func OpenDB() {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"))
	// fmt.Println(connString)

	mysqlConn, err := gorm.Open(mysql.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	Mysql = MysqlDB{
		DB: mysqlConn,
	}

	err = autoMigrate(mysqlConn)
	if err != nil {
		log.Fatal(err)
	}

}

func autoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.Buku{}, // Sesuaikan dengan model yang diinginkan
	)

	if err != nil {
		return err
	}

	return nil
}