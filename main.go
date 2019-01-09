package main

import (
	_ "net/http/pprof"

	"cafapp-returns/gin"
	"cafapp-returns/logger"
)

func main() {
	// // Initalize db
	// logger.Info("Initalizing db...")
	// db, err := gorm.InitDB()
	// if err != nil {
	// 	logger.Fatal("Could not initalize db", err.Error())
	// }

	// //Defer this so that if our application exits, we close the db.
	// defer db.Close()

	// logger.Info("Initalizing Models...")

	// err = gorm.Migrate()
	// if err != nil {
	// 	logger.Info("Could not run object migrations.")
	// }

	logger.Info(`


                         oec :
                        @88888               .d''          .d''
      .         u       8"*88%        u      @8Ne.   .u    @8Ne.   .u
 .udR88N     us888u.    8b.        us888u.   %8888:u@88N   %8888:u@88N
<888'888k .@88 "8888"  u888888> .@88 "8888"   '888I  888.   '888I  888.
9888 'Y"  9888  9888    8888R   9888  9888     888I  888I    888I  888I
9888      9888  9888    8888P   9888  9888     888I  888I    888I  888I
9888      9888  9888    *888>   9888  9888   uW888L  888'  uW888L  888'
?8888u../ 9888  9888    4888    9888  9888  '*88888Nu88P  '*88888Nu88P
 "8888P'  "888*""888"   '888    "888*""888" ~ '88888F'    ~ '88888F'
   "P'     ^Y"   ^Y'     88R     ^Y"   ^Y'     888 ^         888 ^
                         88>                   *8E           *8E
                         48                    '8>           '8>
                         '8                     "             "

`)
	gin.Run()
}
