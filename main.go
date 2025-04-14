package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func main() {
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("andhika", "andhika123", ""),
		Secure: false,
	})

	if err != nil {
		log.Fatalf("Failed to create minio client: %v", err)
	}

	log.Println("MinIO client initialized")

	bucketName := "my-bucket"
	location := "ap-southeast-1"

	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		Region: location,
	})

	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			log.Fatalf("Failed to make minio bucket: %v", err)
		}
	} else {
		log.Printf("Bucket %s created\n", bucketName)
	}

	// upload a file
	imageName := "itachi"
	localFilePath := "img/" + imageName + ".jpg"
	objectName := imageName + ".jpg"

	UploadObject(ctx, minioClient, bucketName, objectName, localFilePath)
	ReadObject(ctx, minioClient, bucketName, objectName)

	// update image pertama ke image kedua, dengan nama image di minio yang sama
	// cuman image yang ganti, namanya enggak
	imageName2 := "itachi2"
	localFilePath2 := "img/" + imageName2 + ".jpg"

	UpdateObject(ctx, minioClient, bucketName, objectName, localFilePath2)
	ReadObject(ctx, minioClient, bucketName, objectName)

	DeleteObject(ctx, minioClient, bucketName, objectName)
	ReadObject(ctx, minioClient, bucketName, objectName)
}

func UploadObject(ctx context.Context, minioClient *minio.Client, bucketName string, objectName string, localFilePath string) {
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, localFilePath, minio.PutObjectOptions{
		ContentType: "image/jpg",
	})

	if err != nil {
		log.Fatalf("Failed to upload object into minio bucket: %v", err)
	}

	fmt.Printf("Uploaded %s (%d bytes)\n", objectName, info.Size)
}

func ReadObject(ctx context.Context, minioClient *minio.Client, bucketName string, objectName string) {
	obj, err := minioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalf("Failed to read %s in minio bucket: %v", objectName, err)
	}

	buf := make([]byte, 10)
	_, err = obj.Read(buf)
	if err != nil {
		log.Fatalf("Failed to read object %s: %v", objectName, err)
	}

	fmt.Printf("Success read %s object from minio bucket\n", objectName)
}

func UpdateObject(ctx context.Context, minioClient *minio.Client, bucketName string, objectName string, localFilePath string) {
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, localFilePath, minio.PutObjectOptions{
		ContentType: "image/jpg",
	})

	if err != nil {
		log.Fatalf("Failed to upload object into minio bucket: %v", err)
	}

	fmt.Printf("Updated to %s From %s (%d bytes)\n", localFilePath, objectName, info.Size)
}

func DeleteObject(ctx context.Context, minioClient *minio.Client, bucketName string, objectName string) {
	err := minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		fmt.Printf("Failed to delete object %s from %s: %v", objectName, bucketName, err)
	}

	fmt.Printf("Object %s is successfully deleted from %s\n", objectName, bucketName)
}
