package main

import (
	pb "awesomeProject1/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStudentServiceClient(conn)

	// 超时设置
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用服务
	r, err := c.GetStudentInfo(ctx, &pb.StudentRequest{Id: "1"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Student: %s", r.GetStudent().GetName())
	log.Printf("Student age: %d", r.GetStudent().GetAge())
	log.Printf("Student class: %s", r.GetStudent().GetClass())
}
