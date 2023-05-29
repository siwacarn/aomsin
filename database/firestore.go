package database

import (
	"context"
	"errors"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/siwacarn/cash-flow-line-bot/models"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type FirestoreDatabase struct {
	ctx        context.Context
	fsClient   *firestore.Client
	collection string
}

func NewFirestoreDatabase(ctx context.Context, collectionName string) (*FirestoreDatabase, error) {
	// 1.) read credential file
	cred := option.WithCredentialsFile("./credential.json")

	// 2.) make firebase app
	app, err := firebase.NewApp(ctx, nil, cred)
	if err != nil {
		return nil, err
	}

	// 3.) firebase app -> firestore
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	cli := &FirestoreDatabase{
		ctx:        ctx,
		fsClient:   client,
		collection: collectionName,
	}

	return cli, nil
}

func (fs *FirestoreDatabase) Create(name string, price int) error {
	_, _, err := fs.fsClient.Collection(fs.collection).Add(fs.ctx, &models.Details{
		Txname:   name,
		Amount:   price,
		Datetime: time.Now(),
	})

	if err != nil {
		return err
	}
	return nil
}

func (fs *FirestoreDatabase) ReadAll() ([]models.Details, error) {
	var model []models.Details

	iter := fs.fsClient.Collection(fs.collection).Documents(fs.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Println(err)
			return nil, err
		}
		// convert document snapshot to model
		// var model models.Details
		var m models.Details
		if err := doc.DataTo(&m); err != nil {
			return nil, err
		}
		model = append(model, m)
	}

	return model, nil
}

func (fs *FirestoreDatabase) ReadOne(id int) (models.Details, error) {
	var datas []models.Details
	iter := fs.fsClient.Collection(fs.collection).Documents(fs.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Println(err)
			return models.Details{}, err
		}
		var m models.Details
		if err := doc.DataTo(&m); err != nil {
			return models.Details{}, err
		}

		// append
		datas = append(datas, m)
	}

	// check len id > datas => err
	if id >= len(datas) && id >= 0 {
		return models.Details{}, errors.New("index out of range")
	}

	return datas[id], nil
}
