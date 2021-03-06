package dto

type InsertTransactionDTO struct {
	ParticipantId uint64  `json:"participant_id" form:"participant_id" binding:"required"`
	EventID       uint64  `json:"event_id" form:"event_id" binding:"required"`
	Amount        float32 `json:"amount" form:"amount" binding:"required"`
	StatusPayment string  `json:"status_payment" form:"status_payment" binding:"required"`
}

type UpdateTransactionDTO struct {
	ID            uint64  `json:"id" form:"id" binding:"required"`
	ParticipantId uint64  `json:"participant_id" form:"participant_id" binding:"required"`
	EventID       uint64  `json:"event_id" form:"event_id" binding:"required"`
	Amount        float32 `json:"amount" form:"amount" binding:"required"`
	StatusPayment string  `json:"status_payment" form:"status_payment" binding:"required"`
}

type BuyEventDTO struct {
	ParticipantId uint64  `json:"participant_id,omitempty" form:"participant_id,omitempty"`
	EventID       uint64  `json:"event_id" form:"event_id" binding:"required"`
	Amount        float32 `json:"amount" form:"amount" binding:"required"`
}

type UploadReceipt struct {
	ParticipantId uint64 `json:"participant_id,omitempty" form:"participant_id,omitempty"`
	ID            uint64 `json:"transaction_id" form:"transaction_id" binding:"required"`
	ImgReceipt    string `json:"receipt" form:"receipt" binding:"required"`
}

type Verify struct {
	TransactionId uint64 `form:"transaction_id" json:"transaction_id" binding:"required"`
	Status        string `form:"status" json:"status" binding:"required"`
}
