package model

import (
	"education/app_server/service"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/sirupsen/logrus"
)

// SetMetadata setMetadata
func SetMetadata(email, A, B, Ci, authority string, keywordIndexByte []byte, serverInfo service.InvokeServer) (response channel.Response, err error) {
	req := channel.Request{
		ChaincodeID: serverInfo.ChaincodeID,
		Fcn:         "setMetadata",
		Args: [][]byte{
			[]byte(email),
			[]byte(authority),
			keywordIndexByte,
			[]byte(A),
			[]byte(B),
			[]byte(Ci),
		},
	}
	response, err = serverInfo.Client.Execute(req)
	if err != nil {
		logrus.Error("setMetadata chaincode func failed,err:", err)
		return
	}
	return
}

// SetKeywordIndex  store keyword to blockchain return index
func SetKeywordIndex(email, keyword, filename string, serverInfo service.InvokeServer) (response channel.Response, err error) {
	req := channel.Request{
		ChaincodeID: serverInfo.ChaincodeID,
		Fcn:         "setKeywordIndex",
		Args: [][]byte{
			[]byte(email),
			[]byte(keyword),
			[]byte(filename),
		},
	}
	response, err = serverInfo.Client.Execute(req)
	if err != nil {
		logrus.Error("set keyword to blockchain return index func failed,err:", err)
		return
	}
	return
}
