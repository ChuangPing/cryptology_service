package sdk

import (
	"education/app_server/sdk/conf"
	"education/app_server/sdk/util"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	cc_name    = "simplecc"
	cc_version = "1.0.0"
)

var (
	//Orgs []*conf.OrgInfo
	Info conf.SdkEnvInfo
)

func InitSDK() (sdk *fabsdk.FabricSDK, err error) {
	// init orgs information
	orgs := []*conf.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/education/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},
	}

	// init sdk env info
	Info = conf.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    os.Getenv("GOPATH") + "/src/education/fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      cc_name,
		//ChaincodePath:    os.Getenv("GOPATH") + "/src/education/chaincode/",
		ChaincodePath:    os.Getenv("GOPATH") + "/src/education/app_server/chaincode/",
		ChaincodeVersion: cc_version,
	}

	// sdk setup
	sdk, err = util.Setup("config.yaml", &Info)
	if err != nil {
		logrus.Error("util setUp failed,err:", err)
		return
	}

	// create channel and join
	if err = util.CreateAndJoinChannel(&Info); err != nil {
		logrus.Error(">> Create channel and join error:", err)
		return
	}

	// create chaincode lifecycle
	if err = util.CreateCCLifecycle(&Info, 1, false, sdk); err != nil {
		logrus.Error(">> create chaincode lifecycle error: %v", err)
		return
	}
	return
}

// 区块链交互
