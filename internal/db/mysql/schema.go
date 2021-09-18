package mysql

import (
	"time"
)

type Register struct {
	// gorm.Model
	Uid          uint32 `gorm:"column:uid;primary_key;auto_increment:false;not null"`
	Openid       string `gorm:"column:openid;type:varchar(64);index;not null"`
	PlatformType uint32 `gorm:"column:platform_type;not null"`
	CreatedTime  int64  `gorm:"column:created_time;not null"`
}

type Player struct {
	// gorm.Model
	Uid       uint32 `gorm:"column:uid;primary_key;auto_increment:false;not null"`
	PbBlob    []byte `gorm:"column:pb_blob;type:blob;not null"`
	UpdatedAt time.Time
}
