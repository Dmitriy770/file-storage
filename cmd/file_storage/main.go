package main

import (
	"context"
	"net/http"

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

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:password@storage:27017/?retryWrites=true&w=majority"))
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
	router.Post("/files/{fileName}", upload.New(bucket))
	router.Get("/files/{fileName}", get.New(bucket))
	router.Delete("/files/{fileName}", delete.New(bucket))
	router.Put("/files/{fileName}", update.New(bucket))
	router.Get("/files", getall.New(bucket))
	router.Get("/files/{fileName}/info", getinfo.New(bucket))

	server := &http.Server{
		Addr:    "app:80",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
