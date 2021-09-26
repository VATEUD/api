package minio

import (
	"api/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"strings"
)

// Minio represents the S3 session
type Minio struct {
	*s3.S3
}

// File represents the File received from the storage
type File struct {
	Name   string
	Object *s3.GetObjectOutput
}

// Upload uploads the file to the storage
func (minio *Minio) Download(path string) (*File, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(utils.Getenv("MINIO_BUCKET", "")),
		Key:    aws.String(path),
	}
	object, err := minio.GetObject(input)

	if err != nil {
		return nil, err
	}

	return &File{minio.fileName(path), object}, nil
}

func (minio Minio) fileName(path string) string {
	split := strings.Split(path, "/")

	return split[len(split)-1]
}

// New starts a new S3 session with the given details
func New() (*Minio, error) {
	config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(utils.Getenv("MINIO_ACCESS_KEY", ""), utils.Getenv("MINIO_SECRET_KEY", ""), ""),
		Endpoint:         aws.String(utils.Getenv("MINIO_ENDPOINT", "")),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(true),
	}

	s3session, err := session.NewSession(config)

	if err != nil {
		return nil, err
	}

	return &Minio{s3.New(s3session)}, nil
}
