package main

import (
    "context"
    "fmt"
    "io"
    "log"

    "google.golang.org/grpc"
    pb "simpleChat/chat"
)

const (
    address = "localhost:12345"
)

func main() {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewChatClient(conn)

    fmt.Print("Enter username: ")
    var user string
    _, err = fmt.Scanln(&user)
    if err != nil {
        log.Fatalf("could not read user: %v", err)
    }
    ctx := context.Background()
    stream, err := c.SendMessage(ctx)
    if err != nil {
        log.Fatal(err)
    }
    waitc := make(chan struct{})
    go func() {
        for {
            in, err := stream.Recv()
            if err == io.EOF {
                close(waitc)
                return
            }
            if err != nil {
                log.Fatalf("failed to receive a message: %v", err)
            }
            fmt.Println(in.Message)
        }
    }()
    for {
        var message string
        _, err = fmt.Scanln(&message)
        if message == "quit" {
            break
        }
        userMsg := &pb.UserMessage{}
        userMsg.User = user
        userMsg.Message = message
        if err := stream.Send(userMsg); err != nil {
            log.Fatalf("failed to send a message: %v", err)
        }
    }

    stream.CloseSend()
    <-waitc
}