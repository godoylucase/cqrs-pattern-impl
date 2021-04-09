package business

import "go.mongodb.org/mongo-driver/bson/primitive"

type URL string
type UserID string
type HashTag string

type Paragraph struct {
	ID       *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Text     string              `json:"text" bson:"text,omitempty"`
	HashTags []HashTag           `json:"hash_tags" bson:"hash_tags,omitempty"`
}

type BaseArticle struct {
	ID             *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID         UserID              `json:"user_id" bson:"user_id,omitempty"`
	Title          string              `json:"title" bson:"title,omitempty"`
	SourceURL      URL                 `json:"source_url" bson:"source_url,omitempty"`
	GlobalHashTags []HashTag           `json:"global_hash_tags" bson:"global_hash_tags,omitempty"`
	Paragraphs     []Paragraph         `json:"paragraphs" bson:"paragraphs,omitempty"`
}
