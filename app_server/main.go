package main

import (
	v1 "education/app_server/api/v1"
	sdk2 "education/app_server/sdk"
	"education/app_server/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()

	// init sdk
	sdk, err := sdk2.InitSDK()
	if err != nil {
		logrus.Error("init sdk failed, err:", err)
		return
	}

	// init service
	service.InitService(sdk2.Info.ChaincodeID, sdk2.Info.ChannelID, sdk2.Info.Orgs[0], sdk)

	// Register
	router.POST("/register", v1.RegisterHandler)
	// login
	router.POST("/login", v1.Login)
	// aes crypto file
	router.POST("/aes", v1.AesHandler)
	// IPFS get store file info
	router.GET("/findAesFile", v1.GetFileInfo)
	//  metadata and Cm store to blockchain
	router.POST("/storeMetadata", v1.MetadataAndCm)
	err = router.Run(":9001")
	if err != nil {
		logrus.Fatalf("run server failed,err:%v\n", err)
	}
}
