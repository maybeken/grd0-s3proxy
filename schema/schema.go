package schema

import (
	"gorm.io/gorm"
)

type S3ProxyTable struct {
	gorm.Model
	ID           uint
	Bucket       string `gorm:"primaryKey"`
	Key          string `gorm:"primaryKey"`
	RequestedAt  int64
	DownloadedAt int64
}
