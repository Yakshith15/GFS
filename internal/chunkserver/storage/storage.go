package storage

import (
	"errors"
	"sync"
	"os"
	"path/filepath"
	"github.com/google/uuid"
	"io"
	"google.golang.org/grpc"
)

type Storage struct{
	baseDir string;
	mu sync.Mutex
	locks map[uuid.UUID]*sync.Mutex;
}

func NewStorage(baseDir string) ( *Storage, error) {
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &Storage{
		baseDir: baseDir,
		mu: sync.Mutex{},
		locks: make(map[uuid.UUID]*sync.Mutex),
	}, nil
}

func (s *Storage) getLock(chunkID uuid.UUID) *sync.Mutex{
	s.mu.Lock()
	defer s.mu.Unlock()
	lock, ok := s.locks[chunkID]
	if ok {
		return lock
	}

	lock = &sync.Mutex{}
	s.locks[chunkID] = lock
	return lock
}

func (s *Storage) getChunkPath(chunkID uuid.UUID) string {
	return filepath.Join(s.baseDir, chunkID.String())
}


func (s *Storage) WriteChunk(chunkID uuid.UUID,  offset int64, data []byte) error {
    lock := s.getLock(chunkID)
	lock.Lock()
	defer lock.Unlock()

	path := s.getChunkPath(chunkID)
	
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteAt(data, offset)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) ReadChunk(chunkID uuid.UUID, offset int64, length int64) ([]byte, error) {
	lock := s.getLock(chunkID)
	lock.Lock()
	defer lock.Unlock()

	path := s.getChunkPath(chunkID)

	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, os.ErrNotExist
		}
		return nil, err
	}
	defer f.Close()

	data := make([]byte, length)

	n, err := f.ReadAt(data, offset)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return data[:n], nil
}

func (s *Storage) ReplicateChunk(chunkID uuid.UUID, sourceAddress string) error {
	lock := s.getLock(chunkID)
	defer lock.Unlock()

	path := s.getChunkPath(chunkID)

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sourceClient, err := grpc.Dial(sourceAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer sourceClient.Close()

}