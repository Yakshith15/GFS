package chunkserver

import (
	"github.com/Yakshith15/GFS/internal/chunkserver/storage"
	"github.com/Yakshith15/GFS/internal/chunkserver/server"
	"log"
	"google.golang.org/grpc"
	"net"
	"github.com/Yakshith15/GFS/api"
)

func main() {
	storage, err := storage.NewStorage("gfs/chunks")
	if err != nil {
		log.Fatalf("failed to create storage: %v", err)
	}
	server := server.NewServer(storage)
	grpcServer := grpc.NewServer()
	api.RegisterChunkServerServiceServer(grpcServer, server)
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Chunk server running on port 5001...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}