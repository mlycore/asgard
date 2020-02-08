package main

import (
	"io"
	"time"
)

type Storage interface {
	ReadFile(path string) (io.ReadCloser, error)
	WriteFile(path string, file io.ReadCloser) error
	ListDirectory(path string) ([]Object, error)
	GetObjectSize(path string) int64
	GetObjectKey(path string) string
	DeleteFile(path string) error
	Copy(src, dst string, recursive bool) error
}

type StorageType string

const (
	StorageTypeS3 = "s3"
	StorageTypeFS = "fs"
)

func NewStorage(config StorageConfig) Storage {
	switch config.Type {
	case "s3":
		return NewS3Storage(config)
	case "fs":
		return NewFileSystemStorage(config)
	default:
		panic("Unknown storage type")
	}
}

type Object interface {
	GetKey() string
	GetLastModified() time.Time
	GetSize() int64
	//	StorageClass() string
}
