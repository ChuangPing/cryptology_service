package service

import (
	"education/app_server/sdk/conf"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/sirupsen/logrus"
)

type InvokeServer struct {
	ChaincodeID string
	Client      *channel.Client
}

var (
	Sdk    *fabsdk.FabricSDK
	Server *InvokeServer
)

func InitService(chaincodeID, channelID string, org *conf.OrgInfo, sdk *fabsdk.FabricSDK) {
	handler := new(InvokeServer)
	handler.ChaincodeID = chaincodeID
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org.OrgUser), fabsdk.WithOrg(org.OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		if err != nil {
			logrus.Error("init service failed, err:", err)
			return
		}
	}
	handler.Client = client
	Server = handler
	Sdk = sdk
}
