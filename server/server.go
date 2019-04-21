package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "simpleChat/chat"
	"time"
)

const (
	port = ":12345"
)

type chatServer struct{}

func (c *chatServer) SendMessage(stream pb.Chat_SendMessageServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		now := time.Now()
		message := &pb.ServerMessage{}
		message.Message = fmt.Sprintf("%d:%d:%d\t%s | %s", now.Hour(), now.Second(), now.Second(), in.User, in.Message)
		if err := stream.Send(message); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServer(s, &chatServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
