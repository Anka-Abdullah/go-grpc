package model

import "time"

type Book struct {
	ID                uint32    `gorm:"primaryKey"`
	Title             string    `gorm:"not null" json:"title"`
	Author            string    `gorm:"not null" json:"author"`
	Publisher         string    `gorm:"not null" json:"publisher"`
	Year              int32     `gorm:"not null" json:"year"`
	ISBN              string    `gorm:"not null;unique" json:"isbn"`
	Quantity          int32     `gorm:"not null" json:"quantity"`
	AvailableQuantity int32     `gorm:"not null" json:"available_quantity"`
	CreatedAt         time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	CreatedBy         uint      `gorm:"not null" json:"created_by"`
}
