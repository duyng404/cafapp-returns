package gorm

import (
	"cafapp-returns/config"
	"cafapp-returns/logger"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite" // sqlite
)

// BaseModel : basic model. may want to factor it out somewhere
type BaseModel struct {
	ID        uint       `gorm:"primary_key,AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func init() {
}

var (
	// DB : the connection object for the db
	DB *gorm.DB
	// for the test db
	// TestDB *gorm.DB
)

// InitDB : initializes the first db, and exports it to be passed around
func InitDB() (*gorm.DB, error) {
	// open the db connection
	logger.Info("Opening a connection to the db...")
	dbFileName := config.DBFilename
	db, err := gorm.Open("sqlite3", "./data/"+dbFileName)
	if err != nil {
		logger.Info("Couldn't open a connection to the db!", err)
		return nil, err
	}

	// if dev, log every query
	if config.ENV == "dev" {
		db.LogMode(true)
		// db.LogMode(false)
	} else {
		db.LogMode(false)
	}

	// Set our variable to use the connection
	DB = db

	return DB, err
}
