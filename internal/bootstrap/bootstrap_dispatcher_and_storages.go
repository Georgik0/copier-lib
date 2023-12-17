package bootstrap

import (
	"context"
	"strconv"

	"storage-api/internal/logger"

	"github.com/lalamove/konfig"
	"github.com/lalamove/konfig/loader/klfile"
	"github.com/lalamove/konfig/parser/kpyaml"
	"github.com/pkg/errors"
)

var configFiles = []klfile.File{
	{
		Path:   "./config/settings.yaml",
		Parser: kpyaml.Parser,
	},
}

const storagesNumberKey = "storages_number"

func init() {
	konfig.Init(konfig.DefaultConfig())
}

type dispatcherUpdater interface {
	AddStorages(ctx context.Context, newStoragesNumber int)
}

func initDispatcherWatcher(ctx context.Context, storageDispatcher dispatcherUpdater) error {
	konfig.RegisterLoaderWatcher(
		klfile.New(&klfile.Config{
			Files: configFiles,
			Watch: true,
		}),

		func(c konfig.Store) error {
			storagesNumber, err := strconv.Atoi(c.String(storagesNumberKey))

			logger.Warn(ctx, "get new storages number: "+c.String(storagesNumberKey))

			if err != nil {
				return errors.Wrap(err, "[initDispatcherWatcher RegisterLoaderWatcher Atoi]")
			}

			storageDispatcher.AddStorages(ctx, storagesNumber)

			return nil
		},
	)

	err := konfig.LoadWatch()

	return errors.Wrap(err, "[initDispatcherWatcher LoadWatch]")
}
