package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/godoylucase/articles_tags/business"
	"github.com/godoylucase/articles_tags/internal"
	"github.com/godoylucase/articles_tags/internal/db"
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

	if err := ensureIndex(coll, "user_id", "source_url"); err != nil {
		return nil, err
	}

	return &repository{
		mongoDBConn: mdbc,
		coll:        coll,
	}, nil
}

func (r *repository) GetByUserIDAndSourceURL(ctx context.Context, userID string, url string) (*business.BaseArticle, error) {
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

func (r *repository) Get(ctx context.Context, id string) (*business.BaseArticle, error) {
	var ba business.BaseArticle

	objID, _ := primitive.ObjectIDFromHex(id)
	if err := r.coll.FindOne(ctx, bson.D{{"_id", objID}}).Decode(&ba); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &ba, nil
}

func (r *repository) Update(ctx context.Context, id string, ba *business.BaseArticle) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objID}}

	updateResult, err := r.coll.ReplaceOne(ctx, filter, ba)
	if err != nil {
		return fmt.Errorf("error when updating article with ID %v and error %w", id, err)
	}

	if updateResult.ModifiedCount == 0 {
		return fmt.Errorf("error article with ID %v not found: %w", id, internal.ErrResourceNotFound)
	}

	return nil
}

// ensureIndex will create index on collection provided
func ensureIndex(cd *mongo.Collection, indexQuery ...string) error {

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
