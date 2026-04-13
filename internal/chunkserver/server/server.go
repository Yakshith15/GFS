package server

import (
	"google.golang.org/grpc"
	"github.com/Yakshith15/GFS/internal/chunkserver/storage"
	"github.com/Yakshith15/GFS/api"
	"github.com/google/uuid"
	"context"
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

func (s *Server) ReadChunk(ctx context.Context, req *api.ReadChunkRequest) (*api.ReadChunkResponse, error) {
	data, err := s.storage.ReadChunk(uuid.MustParse(req.ChunkId), req.Offset, req.Length)
	if err != nil {
		return &api.ReadChunkResponse{
			Status: api.Status_ERROR,
			ErrorMessage: err.Error(),
			Data: data,
		}, nil
	}
	return &api.ReadChunkResponse{
		Status: api.Status_OK,
		ErrorMessage: "",
		Data: data,
	}, nil
}

func (s *Server)WriteChunk(ctx context.Context, req *api.WriteChunkRequest) (*api.WriteChunkResponse, error) {
	err:= s.storage.WriteChunk(uuid.MustParse(req.ChunkId), req.Offset, req.Data)
	if err != nil {
		return &api.WriteChunkResponse{
			Status: api.Status_ERROR,
			ErrorMessage: err.Error(),
		}, nil
	}
	return &api.WriteChunkResponse{
		Status: api.Status_OK,
		ErrorMessage: "",
	}, nil
}

func (s *Server) ReplicateChunk(ctx context.Context, req *api.ReplicateChunkRequest) (*api.ReplicateChunkResponse, error){
	chunkID, err := uuid.Parse(req.ChunkId)
	if err != nil {
		return &api.ReplicateChunkResponse{
			Status: api.Status_ERROR,
			ErrorMessage: err.Error(),
		}, nil
	}
	conn, err := grpc.Dial(req.SourceAddress, grpc.WithInsecure())
	if err != nil {
		return &api.ReplicateChunkResponse{
			Status: api.Status_ERROR,
			ErrorMessage: err.Error(),
		}, nil
	}
	defer conn.Close()
	client := api.NewChunkServerServiceClient(conn)
    data, err := client.ReadChunk(ctx, &api.ReadChunkRequest{
		ChunkId: chunkID.String(),
		Offset:  0,
		Length:  1<<30,
	})
	if err != nil {
		return &api.ReplicateChunkResponse{
			Status: api.Status_ERROR,
			ErrorMessage: err.Error(),
		}, nil
	}
	err = s.storage.WriteChunk(chunkID, 0, data.Data)
	if err != nil {
		return &api.ReplicateChunkResponse{
			Status: api.Status_ERROR,
			ErrorMessage: "failed to write chunk locally",
		}, nil
	}
	return &api.ReplicateChunkResponse{
		Status: api.Status_OK,
		ErrorMessage: "",
	}, nil
}