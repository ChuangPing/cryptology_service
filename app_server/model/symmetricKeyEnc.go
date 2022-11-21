package model

import (
	"education/app_server/service"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/sirupsen/logrus"
)

// store symmetricKeyEnc to blockchain

func StoreSymmetricKeyEnc(email, C1, C2, C3, C4, C5 string, serverInfo service.InvokeServer) (response channel.Response, err error) {
	req := channel.Request{
		ChaincodeID: serverInfo.ChaincodeID,
		Fcn:         "storeSymmetricKeyEnc",
		Args: [][]byte{
			[]byte(email),
			[]byte(C1),
			[]byte(C2),
			[]byte(C3),
			[]byte(C4),
			[]byte(C5),
		},
	}
	response, err = serverInfo.Client.Execute(req)
	if err != nil {
		logrus.Error("client execute failed,err:", err)
		return
	}
	return
}
