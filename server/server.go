package main

import (
	"io"
	"log"
	"net"
	"time"

	pb "simpleChat/chat"

	"google.golang.org/grpc"
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
		in.Timestamp = time.Now().Unix()
		if err := stream.Send(in); err != nil {
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
