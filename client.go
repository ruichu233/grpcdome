package main

import (
	"context"
	"fmt"
	eclient "go.etcd.io/etcd/client/v3"
	eresolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	v1 "grpcdome/proto/shop/v1"
)

func main() {
	serverName := "order_service"

	etcdclient, err := eclient.New(eclient.Config{
		Endpoints: []string{"127.0.0.1:20000"},
	})
	if err != nil {
		panic(err)
	}
	builder, err := eresolver.NewBuilder(etcdclient)
	if err != nil {
		panic(err)
	}
	// 服务名称
	targer := fmt.Sprintf("etcd:///%s", serverName)

	grpcconn, err := grpc.NewClient(targer, grpc.WithResolvers(builder), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = grpcconn.Close()
	}()

	client := v1.NewOrderManagerServiceClient(grpcconn)
	order, err := client.GetOrder(context.Background(), &v1.GetOrderRequest{Id: "1"})
	if err != nil {
		panic(err)
	}

	fmt.Println(order)

}
