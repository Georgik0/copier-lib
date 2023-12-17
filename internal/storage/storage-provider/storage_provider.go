package storage_provider

import (
	abstract_storage "storage-api/internal/storage/abstract-storage"

	"github.com/pkg/errors"
)

var storageProviderIsNilErr = errors.New("storage provider is nil")

type abstractStorage interface {
	SetData(id abstract_storage.FileID, data abstract_storage.Data) error
	LoadData(id abstract_storage.FileID) ([]abstract_storage.Data, error)
}

type Provider struct {
	abstractStorage abstractStorage
	occupancy       int
}

func BuildOne(abstractStorage abstractStorage) Provider {
	return Provider{
		abstractStorage: abstractStorage,
	}
}

func BuildMany(storagesNumber int) []Provider {
	storageProviders := make([]Provider, storagesNumber)

	for i := range storageProviders {
		storageProviders[i] = Provider{
			abstractStorage: abstract_storage.New(i),
		}
	}

	return storageProviders
}

func (p *Provider) SetData(id abstract_storage.FileID, data abstract_storage.Data) error {
	if p == nil {
		return storageProviderIsNilErr
	}

	err := p.abstractStorage.SetData(id, data)
	if err != nil {
		return errors.Wrap(err, "[Provider.SetData]")
	}

	p.updateOccupancy(len(data.Value))

	return nil
}

func (p *Provider) LoadData(id abstract_storage.FileID) ([]abstract_storage.Data, error) {
	if p == nil {
		return nil, storageProviderIsNilErr
	}

	return p.abstractStorage.LoadData(id)
}

func (p *Provider) IsOvercrowded(avgOccupancy float64) bool {
	if p == nil {
		return true
	}

	return avgOccupancy <= float64(p.occupancy)
}

func (p *Provider) updateOccupancy(additionalOccupancy int) {
	if p == nil {
		return
	}

	p.occupancy += additionalOccupancy
}
