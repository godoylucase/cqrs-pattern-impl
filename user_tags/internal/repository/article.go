package repository

import (
	"context"
	"fmt"

	"github.com/godoylucase/cqrs-pattern-impl/business"
	"github.com/godoylucase/cqrs-pattern-impl/internal/db"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	dbName   = "user_articles"
	collName = "articles"
)

type repository struct {
	mongoDBConn *db.MongoDBConn
	coll        *mongo.Collection
}

func NewArticleRepository(mdbc *db.MongoDBConn) (*repository, error) {
	coll := mdbc.Client.Database(dbName).Collection(collName)

	im := mongo.IndexModel{
		Keys: bsonx.MDoc{
			"source_url": bsonx.String("text"),
		},
	}

	indexName, err := coll.Indexes().CreateOne(context.TODO(), im)
	if err != nil {
		return nil, fmt.Errorf("error when initializing indexes from repository with error: %w", err)
	}
	logrus.Infof("index created with name: %v on collection %v", indexName, coll)

	return &repository{
		mongoDBConn: mdbc,
		coll:        coll,
	}, nil
}

func (r *repository) FindByURL(ctx context.Context, url business.URL) (*business.BaseArticle, error) {
	var ba business.BaseArticle
	if err := r.coll.FindOne(ctx, bson.D{{"source_url", url}}).Decode(&ba); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &ba, nil
}

func (r *repository) Create(ctx context.Context, article *business.BaseArticle) (string, error) {
	inserted, err := r.coll.InsertOne(ctx, article)
	if err != nil {
		return "", err
	}

	return inserted.InsertedID.(primitive.ObjectID).String(), nil
}
