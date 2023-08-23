package grpc

import (
	"cook-robot-controller-go/core"
	pb "cook-robot-controller-go/grpc/command"
	"cook-robot-controller-go/logger"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct {
	host string
	port uint16

	controller *core.Controller
}

func NewGRPCServer(host string, port uint16, controller *core.Controller) *GRPCServer {
	return &GRPCServer{
		host:       host,
		port:       port,
		controller: controller,
	}
}

func (g *GRPCServer) Run() {
	address := fmt.Sprintf("%s:%d", g.host, g.port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		logger.Log.Fatalf("无法监听端口: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterCommandServiceServer(server, &command{controller: g.controller})
	logger.Log.Println("gRPC服务启动")

	if err := server.Serve(listen); err != nil {
		logger.Log.Fatalf("无法启动gRPC服务: %v", err)
	}
}
