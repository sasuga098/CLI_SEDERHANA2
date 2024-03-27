package model

import "gorm.io/gorm"

type Buku struct {
	Model
	ISBN   string `gorm:"not null" json:"isbn"`
	Penulis string `gorm:"not null" json:"penulis"`
	Tahun  uint   `gorm:"not null" json:"tahun"`
	Judul  string `gorm:"not null" json:"judul"`
	Gambar string `gorm:"not null" json:"gambar"`
	Stok   uint   `gorm:"not null" json:"stok"`
}

// Create, GetByID, GetAll, UpdateByID, DeleteByID functions...

func (b *Buku) Create(db *gorm.DB) error {
	err := db.Create(&b).Error
	if err != nil {
		return err
	}
	return nil
}

func (b *Buku) GetByID(db *gorm.DB) (Buku, error) {
	res := Buku{}
	err := db.Where("id = ?", b.ID).First(&res).Error
	if err != nil {
		return Buku{}, err
	}
	return res, nil
}

func (b *Buku) GetAll(db *gorm.DB) ([]Buku, error) {
	var res []Buku
	err := db.Find(&res).Error
	if err != nil {
		return []Buku{}, err
	}
	return res, nil
}

func (b *Buku) UpdateByID(db *gorm.DB) error {
	err := db.Save(&b).Error
	if err != nil {
		return err
	}
	return nil
}

func (b *Buku) DeleteByID(db *gorm.DB) error {
	err := db.Delete(&b).Error
	if err != nil {
		return err
	}
	return nil
}