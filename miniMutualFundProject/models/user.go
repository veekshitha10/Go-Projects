package models

import (
	"encoding/json"
	"errors"
	"os"
)

var (
	ErrInvalidName     = errors.New("invalid name field")
	ErrInvalidPassword = errors.New("invalid password ")
	Channel            chan uint
)

type UserTable struct {
	Id       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Image    string `json:"image"`
	Orders []Order `json:"order ,omitempty" gorm:"foreignKey:UserId"`
}

func (u *UserTable) Validate() error {
	if u.Name == "" {
		return ErrInvalidName
	}
	if u.Password == "" {
		return ErrInvalidPassword
	}
	if u.Password==""{
		return  errors.New("invalid password")
	}

	return nil
}
func (u *UserTable) ToBytes() []byte {
	bytes, _ := json.Marshal(u)
	return bytes
}

var (
	Endpoint  = getenv("MINIO_ENDPOINT", "localhost:9000") // your in-cluster DNS:port
	AccessKey = getenv("MINIO_ACCESS_KEY", "minioadmin")
	SecretKey = getenv("MINIO_SECRET_KEY", "minioadmin")
	UseSSL    = getenv("MINIO_USE_SSL", "false") == "true"
	Region    = getenv("MINIO_REGION", "us-east-1") // not useful but still keep some value
	Bucket    = getenv("MINIO_BUCKET", "uploads")
	//object    = getenv("OBJECT_NAME", "hello.txt")
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}