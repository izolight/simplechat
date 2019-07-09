//go:generate protoc -I ../chat --go_out=plugins=grpc:../chat ../chat/chat.proto

package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "time"

    pb "github.com/izolight/simplechat/chat"

    "google.golang.org/grpc"
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
            if in.User == user {
                continue
            }
            ts := time.Unix(in.Timestamp, 0)
            fmt.Printf("%d:%d:%d\t%s | %s\n", ts.Hour(), ts.Minute(), ts.Second(), in.User, in.Text)
        }
    }()
    loop:
    for {
        var text string
        _, err = fmt.Scanln(&text)
        if text == "quit" {
            close(waitc)
            break loop
        }
        userMsg := &pb.ChatMessage{}
        userMsg.User = user
        userMsg.Text = text
        if err := stream.Send(userMsg); err != nil {
            log.Fatalf("failed to send a message: %v", err)
        }
    }

    stream.CloseSend()
    <-waitc
}
