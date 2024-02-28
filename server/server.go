package main

import (
	"awesomeProject1/config"
	pb "awesomeProject1/proto"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
)

type server struct {
	pb.UnimplementedStudentServiceServer
	db *gorm.DB
}

type Student struct {
	ID    string `gorm:"primaryKey"`
	Name  string
	Age   int32
	Class string
}

func NewServer(db *gorm.DB) *server {
	return &server{db: db}
}

func (s *server) GetStudentInfo(ctx context.Context, in *pb.StudentRequest) (*pb.StudentResponse, error) {
	var student Student
	result := s.db.First(&student, "id = ?", in.Id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pb.StudentResponse{Student: &pb.Student{Id: student.ID, Name: student.Name, Age: student.Age, Class: student.Class}}, nil
}

func startGRPCServer(db *gorm.DB) (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStudentServiceServer(s, NewServer(db))

	go func() {
		log.Printf("Starting gRPC server on :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return s, lis
}

func startGRPCGateway() {
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:50051",
		grpc.WithBlock(),
		grpc.WithInsecure(), // Replace with secure connection for production
	)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterStudentServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	gwServer := &http.Server{
		Addr:    ":8080",
		Handler: gwmux,
	}

	log.Printf("Starting gRPC-Gateway on :8080")
	log.Fatal(gwServer.ListenAndServe())
}

func initializeDB() *gorm.DB {
	cfg, err := config.NewConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.Charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&Student{})

	return db
}

func main() {
	db := initializeDB()
	_, lis := startGRPCServer(db)
	defer lis.Close()

	startGRPCGateway()
}
