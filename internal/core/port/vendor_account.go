package port

import "time"

type VendorAccount struct {
	ID         int       `json:"id"`
	VendorType string    `json:"vendor_type" binding:"required"`
	Username   string    `json:"username" binding:"required"`
	Password   string    `json:"password" binding:"required"`
	AppID      *string   `json:"app_id"`
	AppSecret  *string   `json:"app_secret"`
	Token      *string   `json:"token"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
