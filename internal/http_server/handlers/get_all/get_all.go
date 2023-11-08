package getall

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

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

		files := make([]fileDescription, 0)
		for _, file := range foundFiles {
			files = append(files, fileDescription{
				Name:        file.Name,
				Length:      file.Length,
				ContentType: file.Metadata.ContentType,
			})
		}

		response, err := json.Marshal(files)
		if err != nil {
			slog.Error("get all files", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(response)
	}
}
