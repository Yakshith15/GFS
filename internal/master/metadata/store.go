package metadata

import (
	"sync"

	"github.com/Yakshith15/GFS/api"
	"github.com/google/uuid"
)

type MetadataStore struct {
	fileChunks     map[string][]uuid.UUID              // file_path -> ordered chunk_ids
	chunkLocations map[uuid.UUID][]*api.ChunkServer    // chunk_id -> chunkservers
	rw             sync.RWMutex
}

func NewMetadataStore() *MetadataStore {
	return &MetadataStore{
		fileChunks:     make(map[string][]uuid.UUID),
		chunkLocations: make(map[uuid.UUID][]*api.ChunkServer),
	}
}


func (s *MetadataStore) CreateFile(filePath string) api.Status {
	s.rw.Lock()
	defer s.rw.Unlock()

	if _, exists := s.fileChunks[filePath]; exists {
		return api.Status_ALREADY_EXISTS
	}

	s.fileChunks[filePath] = []uuid.UUID{}
	return api.Status_OK
}

func (s *MetadataStore) AddChunk(
	filePath string,
	chunkID uuid.UUID,
	chunkservers []*api.ChunkServer,
) api.Status {

	s.rw.Lock()
	defer s.rw.Unlock()

	if _, exists := s.fileChunks[filePath]; !exists {
		return api.Status_FILE_NOT_FOUND
	}
	s.fileChunks[filePath] = append(s.fileChunks[filePath], chunkID)
	s.chunkLocations[chunkID] = chunkservers
	return api.Status_OK
}

func (s *MetadataStore) GetChunk(
	filePath string,
	chunkIndex int,
) (api.Status, uuid.UUID, []*api.ChunkServer) {

	s.rw.RLock()
	defer s.rw.RUnlock()

	chunks, exists := s.fileChunks[filePath]
	if !exists {
		return api.Status_FILE_NOT_FOUND, uuid.Nil, nil
	}

	if chunkIndex < 0 || chunkIndex >= len(chunks) {
		return api.Status_CHUNK_NOT_FOUND, uuid.Nil, nil
	}

	chunkID := chunks[chunkIndex]
	orig := s.chunkLocations[chunkID]

	serversCopy := make([]*api.ChunkServer, len(orig))
	copy(serversCopy, orig)

	return api.Status_OK, chunkID, serversCopy
}