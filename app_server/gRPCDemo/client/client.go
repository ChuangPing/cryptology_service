package main

import (
	"education/app_server/gRPCDemo/proto/message"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

func main() {
	// 1.Dail连接
	conn, err := grpc.Dial(":8090", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	OrderServiceClient := message.NewOrderServiceClient(conn)

	orderRequest := &message.OrderRequest{
		OrderId:   "201907300001",
		TimeStamp: time.Now().Unix(),
	}
	// 调用函数
	orderInfo, err := OrderServiceClient.GetOrderInfo(context.Background(), orderRequest)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("orderInfo", orderInfo)
}
