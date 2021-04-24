package dto

import (
	"time"
)

type InsertEventDTO struct {
	CreatorId         uint64    `json:"creator_id" form:"creator_id" binding:"required"`
	TitleEvent        string    `json:"title_event" form:"title_event" binding:"required"`
	LinkWebinar       string    `json:"link_webinar" form:"link_webinar" binding:"required"`
	Description       string    `json:"description" form:"description" binding:"required"`
	Banner            string    `json:"banner" form:"banner" binding:"required"`
	Price             float32   `json:"price" form:"price" binding:"required"`
	Quantity          int32     `json:"quantity" form:"quantity" binding:"required"`
	Status            string    `json:"status" form:"status" binding:"required"`
	EventStartDate    time.Time `json:"event_start_date" form:"event_start_date" binding:"required"`
	EventEndDate      time.Time `json:"event_end_date" form:"event_end_date" binding:"required"`
	CampaignStartDate time.Time `json:"campaign_start_date" form:"campaign_start_date" binding:"required"`
	CampaignEndDate   time.Time `json:"campaign_end_date" form:"campaign_end_date" binding:"required"`
}

type UpdateEventDTO struct {
	ID                uint64    `json:"id" form:"id" binding:"required"`
	CreatorId         uint64    `json:"creator_id" form:"creator_id" binding:"required"`
	TitleEvent        string    `json:"title_event" form:"title_event" binding:"required"`
	LinkWebinar       string    `json:"link_webinar" form:"link_webinar" binding:"required"`
	Description       string    `json:"description" form:"description" binding:"required"`
	Banner            string    `json:"banner" form:"banner" binding:"required"`
	Price             float32   `json:"price" form:"price" binding:"required"`
	Quantity          int32     `json:"quantity" form:"quantity" binding:"required"`
	Status            string    `json:"status" form:"status" binding:"required"`
	EventStartDate    time.Time `json:"event_start_date" form:"event_start_date" binding:"required"`
	EventEndDate      time.Time `json:"event_end_date" form:"event_end_date" binding:"required"`
	CampaignStartDate time.Time `json:"campaign_start_date" form:"campaign_start_date" binding:"required"`
	CampaignEndDate   time.Time `json:"campaign_end_date" form:"campaign_end_date" binding:"required"`
}
