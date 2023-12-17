package storage_dispatcher

import (
	"context"
	"fmt"
	"sort"

	"storage-api/internal/logger"
	abstract_storage "storage-api/internal/storage/abstract-storage"
	storage_provider "storage-api/internal/storage/storage-provider"

	"github.com/pkg/errors"
)

var dispatcherIsNilErr = errors.New("Dispatcher is nil")

type Dispatcher struct {
	storageProviders  []storage_provider.Provider       // Слайс со всеми хранилищами
	avgOccupancy      float64                           // Средняя заполненность хранилищ
	storageIdsForData map[abstract_storage.FileID][]int // При сохранении файла запоминаем в какие хранилища он будет записан
}

func New(storagesNumber int) *Dispatcher {
	return &Dispatcher{
		storageProviders:  storage_provider.BuildMany(storagesNumber),
		storageIdsForData: make(map[abstract_storage.FileID][]int),
	}
}

/*
При обновлении конфига будет вызываться этот метод и будут добавлены новые хранилища,
если новое значение больше их текущего количества
После добавления происходит пересчет средней заполненности хранилищ
*/
func (d *Dispatcher) AddStorages(
	ctx context.Context,
	newStoragesNumber int,
) {
	if d == nil {
		logger.Warn(ctx, "Dispatcher is nil")

		return
	}

	providersNumber := len(d.storageProviders)

	additionalStoragesNumber := newStoragesNumber - providersNumber
	if additionalStoragesNumber <= 0 {
		logger.Warn(ctx, "new number of storages less or equal than current number")

		return
	}

	for i := 1; i <= additionalStoragesNumber; i++ {
		d.storageProviders = append(
			d.storageProviders,
			storage_provider.BuildOne(
				abstract_storage.New(providersNumber+i),
			))
	}

	d.avgOccupancy = updateAvgOccupancy(
		providersNumber,
		additionalStoragesNumber,
		d.avgOccupancy,
	)

	warnText := fmt.Sprintf(
		"[Dispatcher.AddStorages] number of storages was successful update; len(storageProviders) = %v; avgOccupancy = %v",
		len(d.storageProviders),
		d.avgOccupancy,
	)

	logger.Warn(ctx, warnText)
}

/*
Для каждой части файла ищется первое доступное хранилище. После записи в хранилище запоминаем ID этого хранилища
*/
func (d *Dispatcher) SetData(
	fileID abstract_storage.FileID, // Идентификатор файла
	dataArr []abstract_storage.Data, // Слайс частей файла. Файл разбили на N частей, где каждая часть - это []byte
) error {
	if d == nil {
		return dispatcherIsNilErr
	}

	// В идеальном варианте каждая запись должна происходить в отдельной горутине,
	// Все записи должны быть обёрнуты в транзакцию, чтобы иметь возможность сделать роллбэк.
	for _, data := range dataArr {
		storageProvider, storageID := d.getFirstAvailableStorageProvider()

		err := storageProvider.SetData(fileID, data)
		if err != nil {
			return errors.Wrap(err, "[Dispatcher.SetData]")
		}

		if len(d.storageIdsForData[fileID]) > 0 {
			continue
		}

		d.storageIdsForData[fileID] = append(d.storageIdsForData[fileID], storageID)
	}

	return nil
}

/*
По идентификатору файла получаем ID хранилищ.
Достаем части файла из каждого хранилища (в некоторых хранилищах может лежать несколько частей файла)
Сортируем части в нужном порядке
*/
func (d *Dispatcher) LoadData(
	fileID abstract_storage.FileID,
) ([]byte, error) {
	if d == nil {
		return nil, dispatcherIsNilErr
	}

	storageIDs := d.storageIdsForData[fileID]

	resultData := make([]byte, 0, len(storageIDs))
	gottenData := make([]abstract_storage.Data, 0, len(storageIDs))

	// В идеале чтение должно выполняться параллельно
	for id := range storageIDs {
		data, err := d.storageProviders[id].LoadData(fileID)
		if err != nil {
			return nil, errors.New("[Dispatcher.LoadData]")
		}

		gottenData = append(gottenData, data...)
	}

	sort.Slice(gottenData, func(i, j int) bool {
		return gottenData[i].Order < gottenData[j].Order
	})

	for i := range gottenData {
		resultData = append(resultData, gottenData[i].Value...)
	}

	return resultData, nil
}

func (d *Dispatcher) getFirstAvailableStorageProvider() (storage_provider.Provider, int) {
	for id, s := range d.storageProviders {
		if s.IsOvercrowded(d.avgOccupancy) {
			continue
		}

		return s, id
	}

	const firstProviderID = 0

	return d.storageProviders[firstProviderID], firstProviderID
}

/*
n - текущее количество хранилищ
avg - ткущая средняя заполненность

avg = (s_1 + s_2 +...+ s_n) / n

После добавления k новых хранилищ
newAvg = (s_1 + s_2 +...+ s_n) + (s_n+1 + s_n+2 + ... + s_n+k) / (n + k)

При условии, что новые хранилища будут изначально пустыми => newAvg = avg * n / (n + k)
*/
func updateAvgOccupancy(
	previousStoragesNumber int,
	additionalStoragesNumber int,
	currentAvgOccupancy float64,
) float64 {
	return currentAvgOccupancy * float64(previousStoragesNumber) /
		float64(previousStoragesNumber+additionalStoragesNumber)
}
