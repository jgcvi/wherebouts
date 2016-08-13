package wherebouts

import (
	//"errors"
	"log"
	//	"os"

	"github.com/gorilla/sessions"
	"gopkg.in/mgo.v2"

	"golang.org/x/net/context"

	//	"google.golang.org/cloud"
	"google.golang.org/cloud/datastore"
	"google.golang.org/cloud/storage"
)

var (
	DB VendorDatabase

	StorageBucket     *storage.BucketHandle
	StorageBucketName string

	SessionStore sessions.Store

	_ mgo.Session
)

func init() {
	var err error
	DB = newMemoryDB()

	// to use Cloud datastore
	DB, err = configureDatastoreDB("wherebouts")

	if err != nil {
		log.Fatal(err)
	}

	// to configure cloud storage
	StorageBucketName = "wherebouts"
	StorageBucket, err = configureStorage(StorageBucketName)

	if err != nil {
		log.Fatal(err)
	}
}

func configureDatastoreDB(projectID string) (VendorDatabase, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return newDatastoreDB(client)
}

func configureStorage(bucketID string) (*storage.BucketHandle, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, err
	}

	return client.Bucket(bucketID), nil
}
