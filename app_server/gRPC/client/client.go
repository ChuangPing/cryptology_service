package gRPCClient

import (
	"context"
	"education/app_server/gRPC/proto/keywordEnc"
	"education/app_server/gRPC/proto/symmetricKeyEnc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//func main() {
//	conn, err := grpc.Dial("222.198.39.39:8011", grpc.WithTransportCredentials(insecure.NewCredentials()))
//	//conn, err := grpc.Dial("222.198.39.39:8011", grpc.WithInsecure())
//	if err != nil {
//		logrus.Fatal("connect failed,err:", err)
//	}
//	defer conn.Close()
//
//	helloServiceClient := hello.NewSayHelloClient(conn)
//
//	request := &hello.Request{
//		Name: "chuanping1",
//	}
//
//	//
//	response, err := helloServiceClient.Hello(context.Background(), request)
//	if err != nil {
//		logrus.Errorf("service failed,err:%v", err)
//	}
//	fmt.Printf(" response:%v\n", response)
//}

// SymmetricKeyEnc client
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
	logrus.Infof("symmetricKeyEnc: C1:%s\t,C2:%s\t,C3:%s\t,C4:%s\t,C5:%s\t", response.C1, response.C2, response.C3, response.C4, response.C5)
	return response.C1, response.C2, response.C3, response.C4, response.C5, nil
}

// KeywordEnc keyword Enc client
func KeywordEnc(keyword string) (A, B, Ci string, err error) {
	conn, err := grpc.Dial("222.198.39.39:8013", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Error("connect keywordEnc server failed,err:", err)
		return
	}
	defer conn.Close()

	keywordClient := keywordEnc.NewKeywordEncServerClient(conn)

	request := keywordEnc.Request{
		Keyword: keyword,
	}

	response, err := keywordClient.KeywordEnc(context.Background(), &request)
	if err != nil {
		logrus.Error("KeywordEnc service failed, err:", err)
		return
	}

	logrus.Infof("success recive keywordEnc A:%s,B:%s,Ci:%s\n", response.A, response.B, response.Ci)
	return response.A, response.B, response.Ci, nil
}
