package models

import "time"

// ImageInterface - общий интерфейс для всех типов изображений
type ImageInterface interface {
	GetID() int
	GetEntityType() string
	GetEntityID() int
	GetFilePath() string
	GetFileName() string
	GetFileSize() int
	GetContentType() string
	GetIsMain() bool
	GetStorageType() string
	GetStorageBucket() string
	GetPublicURL() string
	GetImageURL() string
	GetThumbnailURL() string
	GetDisplayOrder() int
	GetCreatedAt() time.Time
	IsMainImage() bool
	
	SetID(int)
	SetEntityID(int)
	SetFilePath(string)
	SetFileName(string)
	SetFileSize(int)
	SetContentType(string)
	SetIsMain(bool)
	SetStorageType(string)
	SetStorageBucket(string)
	SetPublicURL(string)
	SetImageURL(string)
	SetThumbnailURL(string)
	SetDisplayOrder(int)
	SetCreatedAt(time.Time)
	SetMainImage(bool)
}