package gRPCClient

import (
	"education/app_server/gRPC/proto/hello"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

import (
	"context"
)

func main() {
	conn, err := grpc.Dial("222.198.39.39:8011", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//conn, err := grpc.Dial("222.198.39.39:8011", grpc.WithInsecure())
	if err != nil {
		logrus.Fatal("connect failed,err:", err)
	}
	defer conn.Close()

	helloServiceClient := hello.NewSayHelloClient(conn)

	request := &hello.Request{
		Name: "chuanping1",
	}

	//
	response, err := helloServiceClient.Hello(context.Background(), request)
	if err != nil {
		logrus.Errorf("service failed,err:%v", err)
	}
	fmt.Printf(" response:%v\n", response)
}
