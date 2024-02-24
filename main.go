package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	awsService := NewAwsService()

	// err := awsService.UploadFile("t-ray", "imgs/dice-1.png", "dice-1.png")
	// if err != nil {
	// 	log.Fatalf("update dice image, %v", err)
	// }

	if err := awsService.DownloadFile("t-ray", "imgs/dice-1.png", "dice-1-download.png"); err != nil {
		log.Fatalf("download dice image, %v", err)
	}
}

type AwsService struct {
	S3Client *s3.Client
}

func NewAwsService() *AwsService {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	return &AwsService{
		S3Client: s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UseAccelerate = false
		}),
	}
}

func (s *AwsService) UploadFile(bucketName string, objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
	} else {
		defer file.Close()
		_, err = s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
				fileName, bucketName, objectKey, err)
		}
	}
	return err
}

func (s AwsService) DownloadFile(bucketName string, objectKey string, fileName string) error {
	result, err := s.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return err
}
