package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Storage struct {
	Config     StorageConfig
	S3Client   *s3.S3
	S3Uploader *s3manager.Uploader
}

func NewS3Storage(config StorageConfig) *S3Storage {
	fmt.Printf("%+v\n", config)
	sess := session.Must(session.NewSession())
	creds := credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, "")
	awsOptions := &aws.Config{
		Credentials: creds,
		Region:      aws.String(endpoints.ApNortheast1RegionID),
	}
	svc := s3.New(sess, awsOptions)

	uploader := s3manager.NewUploaderWithClient(svc, func(u *s3manager.Uploader) {
		u.PartSize = 8 * 1024 * 1024
	})

	return &S3Storage{
		Config:     config,
		S3Client:   svc,
		S3Uploader: uploader,
	}
}

func (s *S3Storage) ReadFile(path string) (io.ReadCloser, error) {
	ctx := context.Background()
	result, err := s.S3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Config.BucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (s *S3Storage) WriteFile(path string, file io.ReadCloser) error {
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(s.Config.BucketName),
		Key:    aws.String(path),
		Body:   file,
	}

	_, err := s.S3Uploader.Upload(upParams)
	return err
}

var _ Object = &S3Object{}

type S3Object struct {
	// object *s3.Object
	Key          *string    `json:"key,omitempty"`
	LastModified *time.Time `json:"lastModified,omitempty"`
	Size         *int64     `json:"size,omitempty"`
	// StorageClass *string `json:"storageClass,omitempty"`
}

func (obj *S3Object) GetKey() string {
	return aws.StringValue(obj.Key)
}

func (obj *S3Object) GetLastModified() time.Time {
	return aws.TimeValue(obj.LastModified)
}

func (obj *S3Object) GetSize() int64 {
	return aws.Int64Value(obj.Size)
}

func (s *S3Storage) ListDirectory(path string) ([]Object, error) {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(s.Config.BucketName),
		Prefix: aws.String(strings.TrimPrefix(path, "/")),
	}

	result, err := s.S3Client.ListObjects(params)
	if err != nil {
		return nil, err
	}

	list := make([]Object, len(result.Contents))
	for i := 0; i < len(result.Contents); i++ {
		realkey := strings.TrimPrefix(aws.StringValue(result.Contents[i].Key), strings.TrimPrefix(path, "/"))
		if strings.EqualFold(realkey, "") {
			realkey = "."
		}
		/*
			if strings.EqualFold(strings.Trim(aws.StringValue(result.Contents[i].Key), "/"), strings.Trim(path, "/")) {
				realkey = "./"
			}
		*/
		list[i] = &S3Object{
			//object: result.Contents[i],
			Key:          aws.String(realkey),
			LastModified: result.Contents[i].LastModified,
			Size:         result.Contents[i].Size,
		}
	}

	return list, nil
}

func (s *S3Storage) GetObjectSize(key string) int64 {
	params := &s3.GetObjectInput{
		Bucket: aws.String(s.Config.BucketName),
		Key:    aws.String(key),
	}

	output, err := s.S3Client.GetObject(params)
	if err != nil {
		return 0
	}

	return aws.Int64Value(output.ContentLength)
}

func (s *S3Storage) GetObjectKey(key string) string {
	params := &s3.GetObjectInput{
		Bucket: aws.String(s.Config.BucketName),
		Key:    aws.String(key),
	}

	_, err := s.S3Client.GetObject(params)
	if err != nil {
		return ""
	}

	return key
}

func (s *S3Storage)DeleteFile(file string) error {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String( s.Config.BucketName),
		Key:    aws.String(file),
	}
	_, err := s.S3Client.DeleteObject(params)
	return err
}

func (s *S3Storage)Copy(src, dst string, recursive bool) error {
	// judge if it is a directory
	if strings.EqualFold(s.GetObjectKey(src), "") {
		return errors.New("Key not found")
	}

	if strings.HasSuffix(src, "/") {
		recursive = true
	}

	if recursive {
		logrus.Infof("src: %s", src)
		filelist, err := s.ListDirectory(src)
		if err != nil {
			return err
		}

		for _, file := range filelist {
			path := file.GetKey()
			if strings.EqualFold(path, ".") {
				continue
			}
			logrus.Infof("directory file path: %s", path)
			// readfile
			ctx := context.Background()
			result, err := s.S3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
				Bucket: aws.String(s.Config.BucketName),
				Key:    aws.String(fmt.Sprintf("%s%s", src, path)),
			})
			if err != nil {
				logrus.Errorf("directory file readfile error: err=%s, key=%s", err, fmt.Sprintf("%s/%s", src, path))
				return err
			}

			// changepath and writefile
			upParams := &s3manager.UploadInput{
				Bucket: aws.String(s.Config.BucketName),
				Key:    aws.String(fmt.Sprintf("%s%s", dst, path)),
				Body:   result.Body,
			}

			_, err = s.S3Uploader.Upload(upParams)
			if err != nil {
				logrus.Errorf("directory file writefile error: err=%s, key=%s", err, fmt.Sprintf("%s/%s", dst, path))
				return err
			}
		}

	} else {
		// readfile
		ctx := context.Background()
		result, err := s.S3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.Config.BucketName),
			Key:    aws.String(src),
		})
		if err != nil {
			logrus.Errorf("single file readfile error: err=%s", err)
			return err
		}

		// changepath and writefile
		upParams := &s3manager.UploadInput{
			Bucket: aws.String(s.Config.BucketName),
			Key:    aws.String(dst),
			Body:   result.Body,
		}

		_, err = s.S3Uploader.Upload(upParams)
		if err != nil {
			logrus.Errorf("single file writefile error: err=%s", err)
			return err
		}
	}

	return nil
}
