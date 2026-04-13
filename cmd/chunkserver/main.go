package chunkserver

import (
	"log"
	"net"

	"github.com/Yakshith15/GFS/api"
	"github.com/Yakshith15/GFS/internal/chunkserver/server"
	"github.com/Yakshith15/GFS/internal/chunkserver/storage"
	"google.golang.org/grpc"
)

func main() {
	store, err := storage.NewStorage("./data/chunks")
	if err != nil {
		log.Fatalf("failed to create storage: %v", err)
	}

	cs := server.NewServer(store)

	grpcServer := grpc.NewServer()
	api.RegisterChunkServerServiceServer(grpcServer, cs)

	port := ":5001"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Chunk server running on %s...\n", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}