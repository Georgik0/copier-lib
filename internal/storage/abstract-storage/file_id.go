package abstract_storage

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
)

type FileID string

func ConvFileNameToFileID(strFileID string) (FileID, error) {
	checkSum := md5.Sum([]byte(strFileID))
	key := base64.StdEncoding.EncodeToString(checkSum[:])
	if key == "" {
		return "", fmt.Errorf("[ConvFileNameToFileID] key is empty")
	}

	return FileID(key), nil
}
