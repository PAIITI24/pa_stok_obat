package model

import (
	"time"
)

type StokMasuk struct {
	ID          int       `json:"id" gorm:"primary_key;auto_increment"`
	StokMasuk   uint      `json:"stok_masuk"`
	ExpiredDate time.Time `json:"expired_date"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime" db:"updated_at"`
	ObatID      int       `json:"obat-id"`
	Obat        Obat      `json:"obat,omitempty" gorm:"foreignKey:ObatID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type StokKeluar struct {
	ID         int       `json:"id" gorm:"primary_key;auto_increment"`
	StokKeluar uint      `json:"stok_keluar"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime" db:"updated_at"`
	ObatID     int       `json:"obat-id"`
	Obat       Obat      `json:"obat,omitempty" gorm:"foreignKey:ObatID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
