package business

import "go.mongodb.org/mongo-driver/bson/primitive"

type Paragraph struct {
	Text     string   `json:"text" bson:"text"`
	HashTags []string `json:"hash_tags" bson:"hash_tags"`
}

type BaseArticle struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID         string             `json:"user_id" bson:"user_id"`
	Title          string             `json:"title" bson:"title"`
	SourceURL      string             `json:"source_url" bson:"source_url"`
	GlobalHashTags []string           `json:"global_hash_tags" bson:"global_hash_tags"`
	Paragraphs     []Paragraph        `json:"paragraphs" bson:"paragraphs"`
}

func (ba *BaseArticle) ToDTO() *ArticleDTO {
	id := ba.ID.Hex()
	return &ArticleDTO{
		ID:             id,
		UserID:         ba.UserID,
		Title:          ba.Title,
		SourceURL:      ba.SourceURL,
		GlobalHashTags: ba.GlobalHashTags,
		Paragraphs:     ba.Paragraphs,
	}
}
