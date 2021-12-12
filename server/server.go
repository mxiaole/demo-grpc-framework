package server

import (
	"context"
	"kratos-demo/proto"
)

type GreetServe struct {
	proto.UnimplementedGreeterServer
}

func (g *GreetServe) SayHello(context.Context, *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "hello world"}, nil
}
