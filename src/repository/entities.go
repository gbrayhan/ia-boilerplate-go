package repository

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

type CieVersionType string

const (
	CIE10 CieVersionType = "CIE-10"
	CIE11 CieVersionType = "CIE-11"
)

var ValidCieVersions = []string{
	CIE10.String(),
	CIE11.String(),
}

// check if the cieVersion is valid
func (c CieVersionType) IsValid() bool {
	return c == CIE10 || c == CIE11
}

// return string of the cieVersion
func (c CieVersionType) String() string {
	return string(c)
}

type ICDCie struct {
	ID           int            `gorm:"primaryKey" json:"id"`
	CieVersion   CieVersionType `gorm:"type:varchar(20)" json:"cieVersion"`
	Code         string         `gorm:"type:varchar(20);unique" json:"code"`
	Description  string         `gorm:"type:varchar(255)" json:"description"`
	ChapterNo    string         `gorm:"type:varchar(10)" json:"chapterNo"`
	ChapterTitle string         `gorm:"type:varchar(255)" json:"chapterTitle"`
}

type MedicineType string

const (
	MedicineTypeInjection MedicineType = "injection"
	MedicineTypeTablet    MedicineType = "tablet"
	MedicineTypeCapsule   MedicineType = "capsule"
)

var ValidMedicineTypes = []string{
	MedicineTypeInjection.String(),
	MedicineTypeTablet.String(),
	MedicineTypeCapsule.String(),
}

// check if the medicineType is valid
func (m MedicineType) IsValid() bool {
	return m == MedicineTypeInjection || m == MedicineTypeTablet || m == MedicineTypeCapsule
}

// return string of the medicineType
func (m MedicineType) String() string {
	return string(m)
}

type UnitType string

const (
	UnitTypeMilliliter UnitType = "ml"
	UnitTypeGram       UnitType = "g"
	UnitTypePiece      UnitType = "piece"
	UnitTypeTablet     UnitType = "tablet"
	UnitTypeCapsule    UnitType = "capsule"
)

var ValidUnitTypes = []string{
	UnitTypeMilliliter.String(),
	UnitTypeGram.String(),
	UnitTypePiece.String(),
	UnitTypeTablet.String(),
	UnitTypeCapsule.String(),
}

// check if the unitType is valid
func (u UnitType) IsValid() bool {
	return u == UnitTypeMilliliter || u == UnitTypeGram || u == UnitTypePiece || u == UnitTypeTablet || u == UnitTypeCapsule
}

// return string of the unitType
func (u UnitType) String() string {
	return string(u)
}

type TemperatureControlType string

const (
	TemperatureControlRoom         TemperatureControlType = "room"
	TemperatureControlRefrigerated TemperatureControlType = "refrigerated"
	TemperatureControlFrozen       TemperatureControlType = "frozen"
)

var ValidTemperatureCtrls = []string{
	TemperatureControlRoom.String(),
	TemperatureControlRefrigerated.String(),
	TemperatureControlFrozen.String(),
}

// check if the temperatureControl is valid
func (t TemperatureControlType) IsValid() bool {
	return t == TemperatureControlRoom || t == TemperatureControlRefrigerated || t == TemperatureControlFrozen
}

// return string of the temperatureControl
func (t TemperatureControlType) String() string {
	return string(t)
}

type Medicine struct {
	ID                 int                    `gorm:"primaryKey" json:"id"`
	EANCode            string                 `gorm:"type:varchar(30);unique" json:"eanCode"`
	Description        string                 `gorm:"type:varchar(150)" json:"description"`
	Type               MedicineType           `gorm:"type:varchar(50)" json:"type"`
	Laboratory         string                 `gorm:"type:varchar(50)" json:"laboratory"`
	IVA                string                 `gorm:"type:varchar(5)" json:"iva"`
	SatKey             string                 `gorm:"type:varchar(50)" json:"satKey"`
	TemperatureControl TemperatureControlType `gorm:"type:varchar(50)" json:"temperatureControl"`
	ActiveIngredient   string                 `gorm:"type:varchar(150)" json:"activeIngredient"`
	CreatedAt          time.Time              `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time              `gorm:"autoUpdateTime" json:"updatedAt"`
	ColdChain          bool                   `json:"coldChain"`
	IsControlled       bool                   `json:"isControlled"`
	IsDeleted          bool                   `gorm:"default:false" json:"isDeleted"`
	UnitQuantity       float64                `json:"unitQuantity"`
	UnitType           UnitType               `gorm:"type:varchar(50)" json:"unitType"`
}
