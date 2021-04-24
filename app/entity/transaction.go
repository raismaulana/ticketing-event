package entity

import (
	"database/sql"
	"time"
)

type Transaction struct {
	ID            uint64       `gorm:"primary_key:auto_increment" json:"id"`
	ParticipantId uint64       `json:"participant_id"`
	EventID       uint64       `gorm:"type: not null" json:"event_id"`
	Amount        float32      `gorm:"type:float(15,2)"`
	StatusPayment string       `gorm:"type:enum('passed', 'failed');default:NULL" json:"status_payment"`
	DeletedAt     sql.NullTime `gorm:"type:timestamp null;default:null" json:"deleted_at"`
	CreatedAt     time.Time    `gorm:"type:timestamp not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time    `gorm:"type:timestamp not null ON UPDATE CURRENT_TIMESTAMP();default:CURRENT_TIMESTAMP" json:"Updated_at"`
}

func (Transaction) TableName() string {
	return "Transaction"
}
