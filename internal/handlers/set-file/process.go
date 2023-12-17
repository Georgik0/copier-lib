package set_file

import (
	"context"
	"fmt"
	"net/http"

	make_response "storage-api/internal/common/make-response"
	"storage-api/internal/logger"
	abstract_storage "storage-api/internal/storage/abstract-storage"
	storage_dispatcher "storage-api/internal/storage/storage-dispatcher"

	"github.com/sirupsen/logrus"
)

func NewProcessor(
	logger *logrus.Logger,
	dispatcher *storage_dispatcher.Dispatcher,
) *setFileProcessor {
	return &setFileProcessor{
		logger:     logger,
		dispatcher: dispatcher,
	}
}

type setFileProcessor struct {
	logger     *logrus.Logger
	dispatcher *storage_dispatcher.Dispatcher
}

func (p *setFileProcessor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := logger.ContextWithLogger(context.Background(), p.logger)

	logger.Warn(ctx, "setFiles ServeHTTP process")

	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		make_response.WithError(w, err)
		return
	}

	fileID, err := abstract_storage.ConvFileNameToFileID(fileHeader.Filename)
	if err != nil {
		make_response.WithError(w, err)
		return
	}

	fileMock := []byte{}

	err = p.dispatcher.SetData(fileID, mockSplittingFile(fileMock))
	if err != nil {
		make_response.WithError(w, err)
		return
	}

	logger.Warn(ctx, fmt.Sprintf("fileID: %v", fileID))
}

func mockSplittingFile(file []byte) []abstract_storage.Data {
	// здесь должна быть логика деления файла на куски

	return []abstract_storage.Data{
		{
			Order: 1,
			Value: nil,
		},
		{
			Order: 2,
			Value: nil,
		},
	}
}
