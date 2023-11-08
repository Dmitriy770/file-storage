package delete

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func New(bucket *gridfs.Bucket) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := chi.URLParam(r, "fileName")

		err := bucket.Delete(fileName)
		if errors.Is(err, gridfs.ErrFileNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			slog.Error("delete file", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
