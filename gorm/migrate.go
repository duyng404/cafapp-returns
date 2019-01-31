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
		logger.Error("Error migrating User table:", err)
		return err
	}

	logger.Info("Migrating GoogleUser Table")
	err = DB.AutoMigrate(&GoogleUser{}).Error
	if err != nil {
		logger.Error("Error migrating GoogleUser table:", err)
		return err
	}

	logger.Info("Migrating Global Var Table")
	err = DB.AutoMigrate(&GlobalVar{}).Error
	if err != nil {
		logger.Error("Error migrating Global Var table:", err)
		return err
	}

	logger.Info("Migrating Product Table")
	err = DB.AutoMigrate(&Product{}).Error
	if err != nil {
		logger.Error("Error migrating Product table:", err)
		return err
	}

	logger.Info("Migrating Order Row Table")
	err = DB.AutoMigrate(&OrderRow{}).Error
	if err != nil {
		logger.Error("Error migrating Order Row table:", err)
		return err
	}

	logger.Info("Migrating Destination Table")
	err = DB.AutoMigrate(&Destination{}).Error
	if err != nil {
		logger.Error("Error migrating Destination table:", err)
		return err
	}

	logger.Info("Migrating Order Status Code Table")
	err = DB.AutoMigrate(&OrderStatusCode{}).Error
	if err != nil {
		logger.Error("Error migrating Order Status Code table:", err)
		return err
	}

	logger.Info("Migrating Order Status Update Table")
	err = DB.AutoMigrate(&OrderStatusUpdate{}).Error
	if err != nil {
		logger.Error("Error migrating Order Status Update table:", err)
		return err
	}

	logger.Info("Migrating Order Table")
	err = DB.AutoMigrate(&Order{}).Error
	if err != nil {
		logger.Error("Error migrating Order table:", err)
		return err
	}

	logger.Info("Migrating User Socket Token Table")
	err = DB.AutoMigrate(&UserSocketToken{}).Error
	if err != nil {
		logger.Error("Error migrating user socket token table:", err)
		return err
	}

	logger.Info("Migrating Redeemable Code Table")
	err = DB.AutoMigrate(&RedeemableCode{}).Error
	if err != nil {
		logger.Error("Error migrating redeemable code table:", err)
		return err
	}

	initDestinations()
	initOrderStatusCodes()
	initGlobalVar()

	// init data
	var tmp Product
	err = tmp.PopulateByID(1)
	if err != nil {
		logger.Info("Database is empty, generating sample data")
		initData()
	}

	return nil
}
