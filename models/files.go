package models

import (
	"time"
)

type FileInfo struct {
	ID          string    `bson:"_id"`
	FileName    string    `bson:"fileName"`
	FileURL     string    `bson:"fileUrl"`
	FileType    string    `bson:"fileType"`
	UploadedBy  string    `bson:"uploadedBy"`
	Size        int64     `bson:"size"`
	UploadedAt  time.Time  `bson:"uploadedAt"`
}