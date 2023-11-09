package get

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type gridfsFile struct {
	Name     string `bson:"filename"`
	Length   int64  `bson:"length"`
	Metadata struct {
		ContentType string `bson:"contenttype"`
	} `bson:"metadata"`
}

func New(bucket *gridfs.Bucket) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := chi.URLParam(r, "fileName")

		cursor, err := bucket.Find(bson.D{})
		if err != nil {
			slog.Error("get all files", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var foundFiles []gridfsFile
		err = cursor.All(context.TODO(), &foundFiles)
		if err != nil {
			slog.Error("get all files", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, file := range foundFiles {
			if file.Name == fileName {
				w.Header().Set("Content-Type", file.Metadata.ContentType)
				w.WriteHeader(http.StatusOK)
				break
			}
		}

		_, err = bucket.DownloadToStream(fileName, w)
		if errors.Is(err, gridfs.ErrFileNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			slog.Error("get file", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
