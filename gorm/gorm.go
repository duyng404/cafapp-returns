package gorm

import (
	"cafapp-returns/config"
	"cafapp-returns/logger"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql
)

var (
	// DB : the connection object for the db
	DB *gorm.DB

	// ErrRecordNotFound alias for gorm's ErrRecordNotFound
	ErrRecordNotFound = gorm.ErrRecordNotFound
	// ErrIDZero for when id is zero but it should really not
	ErrIDZero = errors.New("id is zero")
)

// InitDB : initializes the first db, and exports it to be passed around
func InitDB() (*gorm.DB, error) {

	var db *gorm.DB
	var err error

	// Create the Url used to Open the db
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUsername, config.DBPassword, config.DBHostname, config.DBPort, config.DBName)
	logger.Info("dburl is:", dbURL)
	logger.Info("Opening a connection to the db...")
	db, err = gorm.Open("mysql", dbURL)
	if err != nil {
		logger.Info("Couldn't open a connection to the db!", err)
		return nil, err
	}

	// set association autoupdate to false. More info here http://gorm.io/docs/associations.html
	db = db.Set("gorm:association_autoupdate", false)

	// gorm's logging is super f-ing annoying like wtf man why
	// db.LogMode(false)
	db.LogMode(true)

	// Set our variable to use the connection
	DB = db

	return DB, nil
}
