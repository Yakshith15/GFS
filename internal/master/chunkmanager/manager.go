package chunkmanager

import (
	"github.com/Yakshith15/GFS/api"
	"github.com/Yakshith15/GFS/internal/master/metadata"
	"github.com/google/uuid"
)

type ChunkManager struct {
	store            *metadata.MetadataStore
	chunkservers     []*api.ChunkServer
	replicationFactor int
}

func NewChunkManager(
	store *metadata.MetadataStore,
	chunkservers []*api.ChunkServer,
	replicationFactor int,
) *ChunkManager {
	return &ChunkManager{
		store:            store,
		chunkservers:     chunkservers,
		replicationFactor: replicationFactor,
	}
}

func (cm *ChunkManager) AllocateChunk(
	filePath string,
) (api.Status, uuid.UUID, []*api.ChunkServer) {
	if len(cm.chunkservers) < cm.replicationFactor {
		return api.Status_ERROR, uuid.Nil, nil
	}
	chunkID := uuid.New()
	selected := cm.chunkservers[:cm.replicationFactor]
	status := cm.store.AddChunk(filePath, chunkID, selected)
	if status != api.Status_OK {
		return status, uuid.Nil, nil
	}
	return api.Status_OK, chunkID, selected
}