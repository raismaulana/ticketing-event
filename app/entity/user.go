package entity

import (
	"database/sql"
	"time"
)

type User struct {
	ID          uint64        `gorm:"primary_key:auto_increment" json:"id"`
	Username    string        `gorm:"type:varchar(255) unique" json:"username"`
	Fullname    string        `gorm:"type:varchar(255)" json:"name"`
	Email       string        `gorm:"type:varchar(255) unique" json:"email"`
	Password    string        `gorm:"type:text" json:"password"`
	Role        string        `gorm:"type:enum('admin', 'creator', 'participant');default:'participant'" json:"role"`
	Token       string        `gorm:"-" json:"token,omitempty"`
	DeletedAt   sql.NullTime  `gorm:"type:timestamp null;default:null" json:"deleted_at"`
	CreatedAt   time.Time     `gorm:"<-:create;type:timestamp not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"type:timestamp not null ON UPDATE CURRENT_TIMESTAMP();default:CURRENT_TIMESTAMP;" json:"Updated_at"`
	Transaction []Transaction `gorm:"foreignKey:ParticipantId;references:ID" json:"transaction,omitempty"`
	Event       []Event       `gorm:"foreignKey:CreatorId;references:ID" json:"event,omitempty"`
}
