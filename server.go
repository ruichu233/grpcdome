package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"grpcdome/etcd"
	"log"
	"time"

	"google.golang.org/grpc"
	v1 "grpcdome/proto/shop/v1"
	"net"
)

func main() {
	grpcPort, gwPort := ":8009", ":8010"
	go func() {
		listen, err := net.Listen("tcp", grpcPort)
		if err != nil {
			panic(err)
		}
		s := grpc.NewServer(
			grpc.UnaryInterceptor(orderUnaryServerInterceptor),
		)

		v1.RegisterOrderManagerServiceServer(s, &OrderManagerService{})
		go etcd.RegisterEndPointToEtcd(context.Background(), "order_service", grpcPort)
		if err := s.Serve(listen); err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "127.0.0.1"+grpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	gwmux := runtime.NewServeMux()
	if err := v1.RegisterOrderManagerServiceHandler(context.Background(), gwmux, conn); err != nil {
		panic(err)
	}

	//if err := http.ListenAndServe(gwPort, gwmux); err != nil {
	//	panic(err)
	//}
	g := gin.Default()
	g.Any("/v1/orders", func(c *gin.Context) {
		gwmux.ServeHTTP(c.Writer, c.Request)
	})
	g.Run(gwPort)
}

type OrderManagerService struct {
	v1.UnimplementedOrderManagerServiceServer
}

func (s *OrderManagerService) GetOrder(ctx context.Context, req *v1.GetOrderRequest) (*v1.GetOrderResponse, error) {

	return &v1.GetOrderResponse{
		Id:          req.Id,
		Items:       []string{"Google", "Baidu"},
		Description: "example",
		Price:       0,
		Destination: "example",
	}, nil

}

func orderUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	s := time.Now()
	// 拦截器逻辑
	m, err := handler(ctx, req)

	log.Printf("Method: %s,req: %s, latency: %s\n",
		info.FullMethod, req, time.Now().Sub(s),
	)
	return m, err
}
