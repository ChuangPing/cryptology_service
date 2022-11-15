package main

import (
	"context"
	"education/app_server/gRPCDemo/proto/message"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"time"
)

type OrderServiceImpl struct {
	*message.UnimplementedOrderServiceServer
}

//GetOrderInfo(context.Context, *OrderRequest) (*OrderInfo, error)
// 具体方法 GetOrderInfo(context.Context, *OrderRequset) (*OrderInfo, error)  -- 参数：上下文ctx、远程调用传递的参数；返回值：调用的返回值、err
func (os *OrderServiceImpl) GetOrderInfo(ctx context.Context, request *message.OrderRequest) (*message.OrderInfo, error) {
	// 模拟订单数据
	orderMap := map[string]message.OrderInfo{
		"201907300001": {OrderId: "201907300001", OrderName: "衣服", OrderStatus: "已付款"},
		"201907300002": {OrderId: "201907300002", OrderName: "零食", OrderStatus: "已付款"},
		"201907300003": {OrderId: "201907300003", OrderName: "食品", OrderStatus: "未付款"},
	}

	// 指针类型的返回值
	var response *message.OrderInfo
	current := time.Now().Unix()
	if request.TimeStamp > current {
		response = &message.OrderInfo{
			OrderId:     "0",
			OrderName:   "",
			OrderStatus: "订单信息异常",
		}
	} else {
		if result, ok := orderMap[request.OrderId]; ok {
			//*response = result   报错 --指针赋值的坑
			response = &result
		} else {
			response = &message.OrderInfo{
				OrderId:     "0",
				OrderName:   "",
				OrderStatus: "没有订单信息",
			}
		}
	}

	return response, nil
}

func main() {
	server := grpc.NewServer()
	message.RegisterOrderServiceServer(server, new(OrderServiceImpl))

	listen, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err.Error())
	}
	err = server.Serve(listen)
	if err != nil {
		logrus.Fatal("server listen failed,err:", err)
	}
}
