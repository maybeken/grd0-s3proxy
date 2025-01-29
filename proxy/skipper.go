package proxy

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"grd0.net/proxy/s3/database"
	"grd0.net/proxy/s3/localfs"
	"grd0.net/proxy/s3/s3client"
	"grd0.net/proxy/s3/schema"
)

func SkipperLogic(c echo.Context) bool {
	path_segments := strings.Split(c.Request().URL.String(), string(filepath.Separator))

	if len(path_segments) > 2 {
		if c.Request().Method == "GET" {
			return GetHandler(c)
		} else if c.Request().Method == "HEAD" {
			return HeadHandler(c)
		}
	}

	return false
}

func HeadHandler(c echo.Context) bool {
	s3_origin_endpoint := os.Getenv("S3_ORIGIN_ENDPOINT")
	s3_access_key := os.Getenv("S3_ACCESS_KEY")
	s3_secret_access_key := os.Getenv("S3_SECRET_ACCESS_KEY")
	localfs_path := os.Getenv("LOCALFS_PATH")

	bucket := c.Param("bucket")
	object := c.Param("object")

	if object == "" {
		return false
	}

	_, s3_err := s3client.HeadFile(s3_origin_endpoint, s3_access_key, s3_secret_access_key, bucket, object)

	if s3_err != nil {
		os.Remove(localfs_path + "/" + bucket + "/" + object)
		database.DeleteRecord(schema.S3ProxyTable{
			Bucket: bucket,
			Key:    object,
		}, []string{"Bucket", "Key"}, bucket, object)
	}

	return false
}

func GetHandler(c echo.Context) bool {
	s3_origin_endpoint := os.Getenv("S3_ORIGIN_ENDPOINT")
	s3_access_key := os.Getenv("S3_ACCESS_KEY")
	s3_secret_access_key := os.Getenv("S3_SECRET_ACCESS_KEY")
	localfs_path := os.Getenv("LOCALFS_PATH")

	bucket := c.Param("bucket")
	object := c.Param("object")

	file_record := schema.S3ProxyTable{
		Bucket: bucket,
		Key:    object,
	}

	file_exist, err := localfs.FileExists(localfs_path+"/"+bucket, object)
	db_err := database.GetRecord(&file_record)
	remote_file_info, s3_err := s3client.HeadFile(s3_origin_endpoint, s3_access_key, s3_secret_access_key, bucket, object)

	if err != nil && db_err != nil {
		panic("Condition check errors")
	} else if file_exist && s3_err != nil {
		// Normal operation does not call this condition as HEAD OBJECT will reject already
		os.Remove(localfs_path + "/" + bucket + "/" + object)
		database.DeleteRecord(schema.S3ProxyTable{
			Bucket: bucket,
			Key:    object,
		}, []string{"Bucket", "Key"}, bucket, object)
	} else if file_exist && file_record.DownloadedAt >= remote_file_info.LastModified.UnixMilli() {
		return true
	} else {
		err := s3client.DownloadFile(s3_origin_endpoint, s3_access_key, s3_secret_access_key, bucket, object, localfs_path, object)

		if err != nil {
			c.Logger().Fatal(err)
		} else {
			database.UpsertRecord(schema.S3ProxyTable{
				Bucket:       bucket,
				Key:          object,
				RequestedAt:  time.Now().UnixMilli(),
				DownloadedAt: time.Now().UnixMilli(),
			})
		}
	}

	return false
}
