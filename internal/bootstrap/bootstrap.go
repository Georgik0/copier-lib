package bootstrap

import (
	"context"

	storage_dispatcher "storage-api/internal/storage/storage-dispatcher"

	"github.com/pkg/errors"
)

func InitServices(ctx context.Context) (*storage_dispatcher.Dispatcher, error) {
	storageDispatcher := storage_dispatcher.New(0)

	err := initDispatcherWatcher(ctx, storageDispatcher)
	if err != nil {
		return nil, errors.Wrap(err, "[InitServices initDispatcherWatcher]")
	}

	return storageDispatcher, nil
}
