package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type GST struct {
	ID                       uuid.UUID `gorm:"type:text;not null;size:32;" binding:"required"`
	GSTIN                    string    `gorm:"size:30;not null;unique" json:"gstin" binding:"required"`
	TradeName                string    `gorm:"size:255;not null;" json:"tradeName" binding:"required"`
	RegisteredAt             time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"registered_at" binding:"required"`
	Address                  string    `gorm:"size:255;not null;" json:"address" binding:"required"`
	OwnerName                string    `gorm:"size:30;not null;" json:"ownerName" binding:"required"`
	GSTR1LastFiledDate       string    `gorm:"size:15;null;" json:"gstr1LastFiledDate"`
	GSTR3BLastFiledDate      string    `gorm:"size:15;null;" json:"gstr3bLastFiledDate"`
	GSTR9LastFiledDate       string    `gorm:"size:15;null;" json:"gstr9LastFiledDate"`
	GSTR1LastFiledTaxPeriod  string    `gorm:"size:15;null;" json:"gstr1LastFiledTaxPeriod"`
	GSTR3BLastFiledTaxPeriod string    `gorm:"size:15;null;" json:"gstr3bLastFiledTaxPeriod"`
	GSTR9LastFiledTaxPeriod  string    `gorm:"size:15;null;" json:"gstr9LastFiledTaxPeriod"`
	GSTR1PendingReturns      string    `gorm:"type:text;"  json:"gstr1PendingReturns"`
	GSTR3BPendingReturns     string    `gorm:"type:text;"  json:"gstr3bPendingReturns"`
	GSTR9PendingReturns      string    `gorm:"type:text;"  json:"gstr9PendingReturns"`
	CreatedAt                time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *GST) Prepare() {
	p.GSTIN = html.EscapeString(strings.TrimSpace(p.GSTIN))
	p.TradeName = html.EscapeString(strings.TrimSpace(p.TradeName))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *GST) Validate() error {

	if p.GSTIN == "" {
		return errors.New("Required GSTIN")
	}
	if p.OwnerName == "" {
		return errors.New("Required OwnerName")
	}
	if p.TradeName == "" {
		return errors.New("Required TradeName")
	}
	return nil
}

func (p *GST) SaveGST(db *gorm.DB) (*GST, error) {
	var err error
	err = db.Debug().Model(&GST{}).Create(&p).Error
	if err != nil {
		return &GST{}, err
	}
	if p.ID != uuid.Nil {
		err = db.Debug().Model(&User{}).Where("GSTIN = ?", p.GSTIN).Take(&p).Error
		if err != nil {
			return &GST{}, err
		}
	}

	return p, nil
}

func (p *GST) FindAllGSTs(db *gorm.DB) (*[]GST, error) {
	var err error
	posts := []GST{}
	err = db.Debug().Model(&GST{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]GST{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i]).Take(&posts[i]).Error
			if err != nil {
				return &[]GST{}, err
			}
		}
	}
	return &posts, nil
}

func (p *GST) FindGSTByID(db *gorm.DB, gstin string) (*GST, error) {
	var err error
	err = db.Debug().Model(&GST{}).Where("gstin = ?", gstin).Take(&p).Error
	if err != nil {
		return &GST{}, err
	}
	if p.ID != uuid.Nil {
		err = db.Debug().Model(&User{}).Where("id = ?", p).Take(&p).Error
		if err != nil {
			return &GST{}, err
		}
	}
	return p, nil
}

func (p *GST) UpdateAGST(db *gorm.DB) (*GST, error) {

	var err error
	// db = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&Post{}).UpdateColumns(
	// 	map[string]interface{}{
	// 		"title":      p.Title,
	// 		"content":    p.Content,
	// 		"updated_at": time.Now(),
	// 	},
	// )
	// err = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	// if err != nil {
	// 	return &Post{}, err
	// }
	// if p.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
	// 	if err != nil {
	// 		return &Post{}, err
	// 	}
	// }
	err = db.Debug().Model(&GST{}).Where("gstin = ?", p.GSTIN).Updates(GST{GSTIN: p.GSTIN, Address: p.Address, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &GST{}, err
	}
	if p.ID != uuid.Nil {
		err = db.Debug().Model(&User{}).Where("gstin = ?", p.GSTIN).Take(&p).Error
		if err != nil {
			return &GST{}, err
		}
	}
	return p, nil
}

func (p *GST) DeleteAGST(db *gorm.DB, gstin string) (int64, error) {

	db = db.Debug().Model(&GST{}).Where("gstin = ?", gstin).Take(&GST{}).Delete(&GST{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("GST not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
