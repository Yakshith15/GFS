package master

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Yakshith15/GFS/api"
	"github.com/Yakshith15/GFS/internal/master/chunkmanager"
	"github.com/Yakshith15/GFS/internal/master/metadata"
	"github.com/Yakshith15/GFS/internal/master/server"
)

func main() {
	store := metadata.NewMetadataStore()

	chunkservers := []*api.ChunkServer{
		{Id: "cs1", Address: "localhost:5001"},
		{Id: "cs2", Address: "localhost:5002"},
		{Id: "cs3", Address: "localhost:5003"},
	}

	replicationFactor := 3
	cm := chunkmanager.NewChunkManager(store, chunkservers, replicationFactor)

	grpcServer := grpc.NewServer()

	masterServer := server.NewMasterServer(store, cm)
	api.RegisterMasterServiceServer(grpcServer, masterServer)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Master server running on port 8080...")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}