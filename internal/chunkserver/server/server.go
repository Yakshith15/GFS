package server

import (
	"context"
	"net"
	"google.golang.org/grpc"
	"github.com/Yakshith15/GFS/internal/chunkserver/storage"
	"github.com/Yakshith15/GFS/api"
)

type Server struct {
	api.UnimplementedChunkServerServiceServer

	storage *storage.Storage
	grpcServer *grpc.Server
}

func NewServer(storage *storage.Storage) *Server {
	return &Server{
		storage: storage,
		grpcServer: grpc.NewServer(),
	}
}


