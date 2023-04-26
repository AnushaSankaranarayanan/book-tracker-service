package entity

import (
	"time"
)

const (
	defaultUser = "SYSTEM"
	active      = "true"
)

type Book struct {
	ISBN      string `json:"isbn" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Author    string `json:"author" binding:"required"`
	Genre     string `json:"genre" binding:"required"`
	Status    string `json:"status,omitempty"`
	Bookmark  int    `json:"bookmark,omitempty"`
	Created   int64  `json:"created,omitempty"`
	Updated   int64  `json:"updated,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
	Started   int64  `json:"started,omitempty"`
	Finished  int64  `json:"finished,omitempty"`
	Active    string `json:"active,omitempty"`
}

func (b *Book) SetTrackingDetails() {
	b.Created = time.Now().Unix()
	b.Updated = time.Now().Unix()
	b.CreatedBy = defaultUser
	b.UpdatedBy = defaultUser
	b.Active = active
}
