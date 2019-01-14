package gorm

import (
	"cafapp-returns/logger"
)

//Migrate : Initializes all models in db
func Migrate() error {
	logger.Info("Running object Migrations...")

	// One block for each object
	logger.Info("Migrating User Table")
	err := DB.AutoMigrate(&User{}).Error
	if err != nil {
		logger.Error("Error migrating AppSettings table:", err)
		return err
	}

	logger.Info("Migrating GoogleUser Table")
	err = DB.AutoMigrate(&GoogleUser{}).Error
	if err != nil {
		logger.Error("Error migrating AppSettings table:", err)
		return err
	}

	logger.Info("Migrating Product Table")
	err = DB.AutoMigrate(&Product{}).Error
	if err != nil {
		logger.Error("Error migrating AppSettings table:", err)
		return err
	}

	// init data
	var tmp Product
	err = tmp.PopulateByID(1)
	if err != nil {
		logger.Info("Database is empty, generating sample data")
		initData()
	}

	return nil
}
