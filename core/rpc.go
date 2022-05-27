package core

import (
	"net"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	"github.com/jxcia/go-donkey/core/log"
	"google.golang.org/grpc"
)

type server struct {
	// 必须嵌套
	// 包含了 GreeterServer 接口其它方法实现
	pb.UnimplementedGreeterServer
}

/*
func (g *Garden) rpcListen(name, network, address string, obj interface{}, metadata string) error {
	s := server.NewServer()

	rpcx_logger.SetLogger(log.GetLogger())

	if err := s.RegisterName(name, obj, metadata); err != nil {
		return err
	}
	log.Infof("rpc", "listen on: %s", address)
	if err := s.Serve(network, address); err != nil {
		return err
	}
	return nil
}
*/
func (g *Garden) grpcLinsten(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Error("grpc", "", err)
		return err
	}
	s := grpc.NewServer()
	// 将具体的实现注册到服务端
	pb.RegisterGreeterServer(s, &server{})
	// 阻塞等待
	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve:", "", err)
		return err
	}
	return nil
}

/*
创建 gRPC Server 对象，可以把它理解为 Server 端的抽象对象。
将 GreeterServer（其包含需要被调用的服务端接口）注册到 gRPC Server 的内部注册中心，这样在接收请求时，即可通过内部的“服务发现”发现该服务端接口，并进行逻辑处理。
创建 Listen，监听 TCP 端口。
gRPC Server 开始 lis.Accept，直到 Stop 或 GracefulStop
*/
