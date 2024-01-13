package seed

import (
	"log"

	"github.com/google/uuid"
	"github.com/jaganathanb/dapps-api-go/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	{
		ID:       uuid.New(),
		Nickname: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	{
		ID:       uuid.New(),
		Nickname: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

var gsts = []models.GST{
	{
		ID:        uuid.New(),
		GSTIN:     "GSTIN1",
		Address:   "Hello world 1",
		TradeName: "Trade 1",
		OwnerName: "Owner 1",
	},
	{
		ID:        uuid.New(),
		GSTIN:     "GSTIN2",
		Address:   "Hello world 2",
		TradeName: "Trade 2",
		OwnerName: "Owner 2",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.GST{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.GST{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	/*
		err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("attaching foreign key error: %v", err)
		}
	*/

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

		err = db.Debug().Model(&models.GST{}).Create(&gsts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}
