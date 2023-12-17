package main

import (
	"context"
	"log"
	"net/http"

	"storage-api/internal/bootstrap"
	set_file "storage-api/internal/handlers/set-file"
	"storage-api/internal/logger"

	"github.com/pkg/errors"
)

const errListenAndServe = "ListenAndServe"

func main() {
	ctx := context.Background()

	l := logger.InitLogger()
	ctx = logger.ContextWithLogger(ctx, l)
	logger.Warn(ctx, "service started...")

	dispatcher, err := bootstrap.InitServices(ctx)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	//mux.Handle("/getFile", nil)
	mux.Handle("/setFile", set_file.NewProcessor(
		l,
		dispatcher,
	))

	err = http.ListenAndServe(":9000", mux)
	if err != nil {
		logger.Warn(ctx, errors.Wrap(err, errListenAndServe).Error())
	}
}
