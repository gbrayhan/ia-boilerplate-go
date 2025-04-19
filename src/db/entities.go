package db

import (
	"fmt"
	"time"
)

func MigrateEntitiesGORM() error {
	err := DB.AutoMigrate(&User{}, &RoleUser{}, &DeviceDetails{}, &Medicine{}, &ICDCie{}, &Client{}, &SubClient{}, &ProgramSubclient{})
	if err != nil {
		fmt.Println("Error migrating the database:", err)
		return err
	}
	return nil
}

type RoleUser struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Role        string    `gorm:"unique;not null" json:"role"`
	Description string    `json:"description"`
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type User struct {
	ID           int             `gorm:"primaryKey" json:"id"`
	Username     string          `gorm:"unique;not null" json:"username"`
	FullName     string          `json:"fullName"`
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

type Client struct {
	ID             int         `gorm:"primaryKey" json:"id"`
	Alias          string      `gorm:"size:100" json:"alias"`
	LegalName      string      `gorm:"size:100" json:"legalName"`
	TIN            string      `gorm:"size:15" json:"tin"`
	ContractNumber string      `gorm:"size:100" json:"contractNumber"`
	FiscalAddress  string      `gorm:"size:255" json:"fiscalAddress"`
	CreatedAt      time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time   `gorm:"autoUpdateTime" json:"updatedAt"`
	IsDeleted      int         `gorm:"default:0" json:"isDeleted"`
	SubClients     []SubClient `gorm:"foreignKey:ClientID" json:"subClients"`
}

type SubClient struct {
	ID             int                `gorm:"primaryKey" json:"id"`
	ClientID       int                `json:"clientId"`
	Alias          string             `gorm:"size:100" json:"alias"`
	ContractNumber string             `gorm:"size:100" json:"contractNumber"`
	LegalName      string             `gorm:"size:100" json:"legalName"`
	FiscalAddress  string             `gorm:"size:255" json:"fiscalAddress"`
	TIN            string             `gorm:"size:15" json:"tin"`
	CreatedAt      time.Time          `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time          `gorm:"autoUpdateTime" json:"updatedAt"`
	IsDeleted      int                `gorm:"default:0" json:"isDeleted"`
	Programs       []ProgramSubclient `gorm:"foreignKey:SubclientID" json:"programs"`
}

type ProgramSubclient struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	SubclientID int       `json:"subclientId"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	IsDeleted   int       `gorm:"default:0" json:"isDeleted"`
}

type ICDCie struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	CieVersion   string `gorm:"type:varchar(2)" json:"cieVersion"`
	Code         string `gorm:"type:varchar(20)" json:"code"`
	Description  string `gorm:"type:varchar(255)" json:"description"`
	ChapterNo    string `gorm:"type:varchar(10)" json:"chapterNo"`
	ChapterTitle string `gorm:"type:varchar(255)" json:"chapterTitle"`
}

type Medicine struct {
	ID               int       `gorm:"primaryKey" json:"id"`
	EANCode          string    `gorm:"type:varchar(30)" json:"eanCode"`
	Description      string    `gorm:"type:varchar(150)" json:"description"`
	Type             string    `gorm:"type:varchar(50)" json:"type"`
	Laboratory       string    `gorm:"type:varchar(50)" json:"laboratory"`
	IVA              string    `gorm:"type:varchar(5)" json:"iva"`
	SatKey           string    `gorm:"type:varchar(50)" json:"satKey"`
	ActiveIngredient string    `gorm:"type:varchar(150)" json:"activeIngredient"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	ColdChain        bool      `json:"coldChain"`
	IsControlled     bool      `json:"isControlled"`
	IsDeleted        bool      `gorm:"default:false" json:"isDeleted"`
	UnitQuantity     float64   `json:"unitQuantity"`
	UnitType         string    `gorm:"type:varchar(50)" json:"unitType"`
}
