package s3client

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func DownloadFile(endpoint, accessKeyID, secretAccessKey, bucket, objectKey, path, filename string) error {
	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure:       false, // Set to true if using HTTPS
		BucketLookup: 2,
	})

	if err != nil {
		return fmt.Errorf("failed to create MinIO client: %w", err)
	}

	err = os.MkdirAll(path+"/"+bucket, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Create a file to write the downloaded content
	outFile, err := os.Create(path + "/" + bucket + "/" + filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	decoded_object_key, err := url.QueryUnescape(objectKey)
	// Download the file
	err = minioClient.FGetObject(context.Background(), bucket, decoded_object_key, path+"/"+bucket+"/"+filename, minio.GetObjectOptions{})

	if err != nil {
		os.Remove(path + "/" + bucket + "/" + filename)
		return fmt.Errorf("failed to download file: %w", err)
	}

	return nil
}

func HeadFile(endpoint, accessKeyID, secretAccessKey, bucket, objectKey string) (minio.ObjectInfo, error) {
	if bucket == "" || objectKey == "" {
		panic("Bucket / Object key cannot be empty")
	}

	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure:       false, // Set to true if using HTTPS
		BucketLookup: 2,
	})

	if err != nil {
		panic(err)
	}

	decoded_object_key, err := url.QueryUnescape(objectKey)
	// Head the file
	objectInfo, err := minioClient.StatObject(context.Background(), bucket, decoded_object_key, minio.StatObjectOptions{})

	return objectInfo, err
}
