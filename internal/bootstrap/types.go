package bootstrap

import (
	"context"

	abstract_storage "storage-api/internal/storage/abstract-storage"
)

type storageDispatcherI interface {
	AddStorages(ctx context.Context, newStoragesNumber int)
	SetData(id abstract_storage.FileID, data []byte) error
	LoadData(id abstract_storage.FileID) ([]byte, error)
}
