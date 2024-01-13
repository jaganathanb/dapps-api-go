package controllertests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jaganathanb/dapps-api-go/api/controllers"
	"github.com/jaganathanb/dapps-api-go/api/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var userInstance = models.User{}
var gstInstance = models.GST{}

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())

}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "sqlite3" {
		//DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		testDbName := os.Getenv("TestDbName")
		server.DB, err = gorm.Open(TestDbDriver, testDbName)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
		server.DB.Exec("PRAGMA foreign_keys = ON")
	}
}

func refreshUserTable() error {
	/*
		err := server.DB.DropTableIfExists(&models.User{}).Error
		if err != nil {
			return err
		}
		err = server.DB.AutoMigrate(&models.User{}).Error
		if err != nil {
			return err
		}
	*/
	err := server.DB.DropTableIfExists(&models.GST{}, &models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.GST{}).Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed table(s)")
	return nil
}

func seedOneUser() (models.User, error) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Nickname: "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func seedUsers() ([]models.User, error) {

	var err error
	if err != nil {
		return nil, err
	}
	users := []models.User{
		{
			Nickname: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		{
			Nickname: "Kenny Morris",
			Email:    "kenny@gmail.com",
			Password: "password",
		},
	}
	for i := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}

func refreshUserAndGstTable() error {

	err := server.DB.DropTableIfExists(&models.GST{}, &models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.GST{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOneGst() (models.GST, error) {

	err := refreshUserAndGstTable()
	if err != nil {
		return models.GST{}, err
	}
	user := models.User{
		Nickname: "Sam Phil",
		Email:    "sam@gmail.com",
		Password: "password",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.GST{}, err
	}
	gst := models.GST{
		GSTIN:     "GSTIN1",
		TradeName: "This is the content sam",
		Address:   "Addr 1",
	}
	err = server.DB.Model(&models.GST{}).Create(&gst).Error
	if err != nil {
		return models.GST{}, err
	}
	return gst, nil
}

func seedUsersAndGsts() ([]models.User, []models.GST, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.GST{}, err
	}
	var users = []models.User{
		{
			Nickname: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		{
			Nickname: "Magu Frank",
			Email:    "magu@gmail.com",
			Password: "password",
		},
	}
	var gsts = []models.GST{
		{
			GSTIN:     "GSTIN1",
			TradeName: "Hello world 1",
		},
		{
			GSTIN:     "GSTIN2",
			TradeName: "Hello world 2",
		},
	}

	for i := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

		err = server.DB.Model(&models.GST{}).Create(&gsts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed gsts table: %v", err)
		}
	}

	return users, gsts, nil
}
