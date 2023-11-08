package update

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(bucket *gridfs.Bucket) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		fileName := chi.URLParam(r, "fileName")

		err := bucket.Delete(fileName)
		if errors.Is(err, gridfs.ErrFileNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		uploadOpts := options.GridFSUpload().SetMetadata(bson.D{{Key: "contenttype", Value: r.Header.Get("Content-Type")}})
		err = bucket.UploadFromStreamWithID(fileName, fileName, io.Reader(r.Body), uploadOpts)
		if mongo.IsDuplicateKeyError(err) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		if err != nil {
			slog.Error("update file", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
