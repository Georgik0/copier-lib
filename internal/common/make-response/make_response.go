package make_response

import "net/http"

func WithError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}
