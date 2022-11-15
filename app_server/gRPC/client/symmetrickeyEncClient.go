package gRPCClient

import (
	"context"
	"education/app_server/gRPC/proto/symmetricKeyEnc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//
func SymmetricKeyEnc(symmetricKey string) (C1, C2, C3, C4, C5 string, err error) {
	//func main() {

	conn, err := grpc.Dial("222.198.39.39:8012", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatal("connect symmetricKeyEnc Server failed,err:", err)
		return
	}
	defer conn.Close()
	symmetricKeyEncClient := symmetricKeyEnc.NewSymmetricKeyEncServerClient(conn)
	request := symmetricKeyEnc.Request{
		SymmetricKey: symmetricKey,
	}
	response, err := symmetricKeyEncClient.SymmetricKeyEnc(context.Background(), &request)
	if err != nil {
		logrus.Error("symmetricKeyEnc server failed,err:", err)
		return
	}
	logrus.Info("enc:", response.C1, response.C2, response.C3, response.C4, response.C5)
	return response.C1, response.C2, response.C3, response.C4, response.C5, nil
}
