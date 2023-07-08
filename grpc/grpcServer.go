package grpc

import (
	"cook-robot-controller-go/core"
	pb "cook-robot-controller-go/grpc/commandServer"
	"cook-robot-controller-go/logger"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct {
	controller *core.Controller
}

func NewGRPCServer(controller *core.Controller) *GRPCServer {
	return &GRPCServer{
		controller: controller,
	}
}

func (g *GRPCServer) Run() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log.Fatalf("无法监听端口: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterCommandServiceServer(server, &commandServer{controller: g.controller})
	logger.Log.Println("gRPC服务已启动，等待客户端连接...")

	if err := server.Serve(listen); err != nil {
		logger.Log.Fatalf("无法启动gRPC服务: %v", err)
	}
}
