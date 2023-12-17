package abstract_storage

import (
	"errors"
)

var (
	storageIsNilErr = errors.New("storage is nil")
	cantFoundErr    = errors.New("can't find value for the key")
)

type abstractStorage struct {
	id   int
	data map[FileID][]Data
}

func New(storageID int) *abstractStorage {
	return &abstractStorage{
		id:   storageID,
		data: make(map[FileID][]Data),
	}
}

func (s *abstractStorage) SetData(id FileID, data Data) error {
	if s == nil {
		return storageIsNilErr
	}

	s.data[id] = append(s.data[id], data)

	return nil
}

func (s *abstractStorage) LoadData(id FileID) ([]Data, error) {
	if s == nil {
		return nil, storageIsNilErr
	}

	if data, ok := s.data[id]; ok {
		return data, nil
	}

	return nil, cantFoundErr
}
