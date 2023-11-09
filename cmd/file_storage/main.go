package main

import (
	"context"
	"net/http"
	"os"

	"github.com/Dmitriy770/file-storage/internal/http_server/handlers/delete"
	"github.com/Dmitriy770/file-storage/internal/http_server/handlers/get"
	getall "github.com/Dmitriy770/file-storage/internal/http_server/handlers/get_all"
	getinfo "github.com/Dmitriy770/file-storage/internal/http_server/handlers/get_info"
	"github.com/Dmitriy770/file-storage/internal/http_server/handlers/update"
	"github.com/Dmitriy770/file-storage/internal/http_server/handlers/upload"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ADDRESS         string
	MONGODB_ADDRESS string
)

func init() {
	ADDRESS = os.Getenv("ADDRESS")
	MONGODB_ADDRESS = os.Getenv("MONGODB_ADDRESS")
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGODB_ADDRESS))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := client.Database("filesDB")
	bucket, err := gridfs.NewBucket(db)

	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	router.Route("/files/{fileName}", func(r chi.Router) {
		r.Post("/", upload.New(bucket))
		r.Get("/", get.New(bucket))
		r.Delete("/", delete.New(bucket))
		r.Put("/", update.New(bucket))
		r.Get("/info", getinfo.New(bucket))
	})
	router.Get("/files", getall.New(bucket))

	server := &http.Server{
		Addr:    ADDRESS,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
