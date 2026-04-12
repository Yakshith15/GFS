package server

import (
	"context"

	"github.com/Yakshith15/GFS/api"
	"github.com/Yakshith15/GFS/internal/master/chunkmanager"
	"github.com/Yakshith15/GFS/internal/master/metadata"
)

type MasterServer struct {
	api.UnimplementedMasterServiceServer

	store         *metadata.MetadataStore
	chunkManager  *chunkmanager.ChunkManager
}

func NewMasterServer(
	store *metadata.MetadataStore,
	chunkManager *chunkmanager.ChunkManager,
) *MasterServer {
	return &MasterServer{
		store:        store,
		chunkManager: chunkManager,
	}
}

func (s *MasterServer) CreateFile(
	ctx context.Context,
	req *api.CreateFileRequest,
) (*api.CreateFileResponse, error) {

	status := s.store.CreateFile(req.FilePath)

	return &api.CreateFileResponse{
		Status:       status,
		ErrorMessage: "",
	}, nil
}

func (s *MasterServer) AllocateChunk(
	ctx context.Context,
	req *api.AllocateChunkRequest,
) (*api.AllocateChunkResponse, error) {

	status, chunkID, chunkservers := s.chunkManager.AllocateChunk(req.FilePath)

	if status != api.Status_OK {
		return &api.AllocateChunkResponse{
			Status:       status,
			ErrorMessage: "failed to allocate chunk",
		}, nil
	}

	return &api.AllocateChunkResponse{
		Status:       api.Status_OK,
		ErrorMessage: "",
		ChunkId:      chunkID.String(),
		Chunkservers: chunkservers,
	}, nil
}

func (s *MasterServer) GetChunkLocations(
	ctx context.Context,
	req *api.GetChunkLocationRequest,
) (*api.GetChunkLocationResponse, error) {

	status, chunkID, chunkservers := s.store.GetChunk(
		req.FilePath,
		int(req.ChunkIndex),
	)

	if status != api.Status_OK {
		return &api.GetChunkLocationResponse{
			Status:       status,
			ErrorMessage: "chunk not found",
		}, nil
	}

	return &api.GetChunkLocationResponse{
		Status:       api.Status_OK,
		ErrorMessage: "",
		ChunkId:      chunkID.String(),
		Chunkservers: chunkservers,
	}, nil
}