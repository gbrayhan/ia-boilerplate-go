package repository

// IMPORTANT: in this file you will find the entities that are used in the database and only in the database

import (
	"time"
)

type RoleUser struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type User struct {
	ID           int             `gorm:"primaryKey" json:"id"`
	Username     string          `gorm:"unique;not null" json:"username"`
	FirstName    string          `json:"firstName"`
	LastName     string          `json:"lastName"`
	Email        string          `gorm:"unique;not null" json:"email"`
	HashPassword string          `gorm:"not null" json:"-"`
	JobPosition  string          `json:"jobPosition"`
	RoleID       int             `json:"roleId"`
	Role         RoleUser        `gorm:"foreignKey:RoleID" json:"role"`
	Enabled      bool            `gorm:"default:true" json:"enabled"`
	CreatedAt    time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`
	Devices      []DeviceDetails `gorm:"foreignKey:UserID" json:"devices"`
}

type DeviceDetails struct {
	ID             int       `gorm:"primaryKey" json:"id"`
	UserID         int       `gorm:"not null" json:"userId"`
	IPAddress      string    `gorm:"type:varchar(45);not null" json:"ip_address"`
	UserAgent      string    `json:"user_agent"`
	DeviceType     string    `json:"device_type"`
	Browser        string    `json:"browser"`
	BrowserVersion string    `json:"browser_version"`
	OS             string    `json:"os"`
	Language       string    `json:"language"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
type ICDCie struct {
	ID int `gorm:"primaryKey" json:"id"`
	// CieVersion holds the ICD code version string like "CIE-10" or "CIE-11".
	// The previous size of 2 characters was too small and caused database
	// errors when inserting typical values, so it has been increased.
	CieVersion   string `gorm:"type:varchar(20)" json:"cieVersion"`
	Code         string `gorm:"type:varchar(20)" json:"code"`
	Description  string `gorm:"type:varchar(255)" json:"description"`
	ChapterNo    string `gorm:"type:varchar(10)" json:"chapterNo"`
	ChapterTitle string `gorm:"type:varchar(255)" json:"chapterTitle"`
}

type Medicine struct {
	ID               int       `gorm:"primaryKey" json:"id"`
	EANCode          string    `gorm:"type:varchar(30);unique" json:"eanCode"`
	Description      string    `gorm:"type:varchar(150)" json:"description"`
	Type             string    `gorm:"type:varchar(50)" json:"type"`
	Laboratory       string    `gorm:"type:varchar(50)" json:"laboratory"`
	IVA              string    `gorm:"type:varchar(5)" json:"iva"`
	SatKey           string    `gorm:"type:varchar(50)" json:"satKey"`
	TemperatureCtrl  string    `gorm:"type:varchar(50)" json:"temperatureControl"`
	ActiveIngredient string    `gorm:"type:varchar(150)" json:"activeIngredient"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	ColdChain        bool      `json:"coldChain"`
	IsControlled     bool      `json:"isControlled"`
	IsDeleted        bool      `gorm:"default:false" json:"isDeleted"`
	UnitQuantity     float64   `json:"unitQuantity"`
	UnitType         string    `gorm:"type:varchar(50)" json:"unitType"`
}
