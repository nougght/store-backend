package models

import (
	"time"
)

type Url struct {
	ObjectName string    `db:"object_name" json:"object_name"`
	BucketName string    `db:"bucket_name" json:"bucket_name"`
	Url        string    `db:"url" json:"url"`
	ExpiresAt  time.Time `db:"expires_at" json:"expires_at"`
}



