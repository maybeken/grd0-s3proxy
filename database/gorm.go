package database

import (
	"os"
	"strings"

	"github.com/thoas/go-funk"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"grd0.net/proxy/s3/schema"
)

var db *gorm.DB

func InitDatabase() {
	localfs_path := os.Getenv("LOCALFS_PATH")

	db_instance, err := gorm.Open(sqlite.Open(localfs_path+"/s3_proxy.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db = db_instance

	db.AutoMigrate(&schema.S3ProxyTable{})
}

func UpsertRecord[T any](data T) {
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&data)
}

func DeleteRecord[T any](data T, pk []string, values ...interface{}) {
	whereClauses := funk.Map(pk, func(key string) string {
		return key + " = ?"
	}).([]string)

	db.Where(strings.Join(whereClauses, " AND "), values...).Delete(&data)
}

func GetRecord[T any](data T) *gorm.DB {
	return db.Take(&data)
}
