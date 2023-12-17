package storage_dispatcher

import (
	"fmt"
	"sort"
	"testing"

	abstract_storage "storage-api/internal/storage/abstract-storage"
)

func TestMy(t *testing.T) {
	gottenData := make([]abstract_storage.Data, 0, 10)

	for i := 9; i >= 0; i-- {
		gottenData = append(gottenData, abstract_storage.Data{
			Value: nil,
			Order: i,
		})
	}

	fmt.Println(gottenData)

	sort.Slice(gottenData, func(i, j int) bool {
		return gottenData[i].Order < gottenData[j].Order
	})

	fmt.Println(gottenData)
}
