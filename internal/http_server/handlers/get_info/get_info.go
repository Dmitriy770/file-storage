package getinfo

import (
	"context"
	"encoding/json"
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

type fileDescription struct {
	Name        string `bson:"filename"`
	Length      int64  `bson:"length"`
	ContentType string `bson:"contenttype"`
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
				response, err := json.Marshal(fileDescription{
					Name:        file.Name,
					Length:      file.Length,
					ContentType: file.Metadata.ContentType,
				})
				if err != nil {
					slog.Error("get all files", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.Write(response)
				return
			}
		}
	}
}
