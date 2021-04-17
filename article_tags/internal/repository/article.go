package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/godoylucase/cqrs-pattern-impl/business"
	"github.com/godoylucase/cqrs-pattern-impl/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	dbName   = "articles"
	collName = "articles"
)

type repository struct {
	mongoDBConn *db.MongoDBConn
	coll        *mongo.Collection
}

func NewArticleRepository(mdbc *db.MongoDBConn) (*repository, error) {
	coll := mdbc.Client.Database(dbName).Collection(collName)

	if err := EnsureIndex(coll, "user_id", "source_url"); err != nil {
		return nil, err
	}

	return &repository{
		mongoDBConn: mdbc,
		coll:        coll,
	}, nil
}

func (r *repository) GetByUserIDAndSourceURL(ctx context.Context, userID string, url business.URL) (*business.BaseArticle, error) {
	var ba business.BaseArticle
	if err := r.coll.FindOne(ctx, bson.D{
		{"source_url", url},
		{"user_id", userID},
	}).Decode(&ba); err != nil {
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

	oid := inserted.InsertedID.(primitive.ObjectID)

	article.ID = oid

	return oid.Hex(), nil
}

// EnsureIndex will create index on collection provided
func EnsureIndex(cd *mongo.Collection, indexQuery ...string) error {

	opts := options.CreateIndexes().SetMaxTime(3 * time.Second)

	var index []mongo.IndexModel

	for _, val := range indexQuery {
		temp := mongo.IndexModel{}
		temp.Keys = bsonx.Doc{{Key: val, Value: bsonx.Int32(1)}}
		index = append(index, temp)
	}
	_, err := cd.Indexes().CreateMany(context.Background(), index, opts)
	if err != nil {
		fmt.Errorf("error while executing index query, with error: %w", err.Error())
		return err
	}
	return nil
}
