package dto

type ArticleDTO struct {
	ID             string   `json:"id"`
	UserID         string   `json:"user_id"`
	Title          string   `json:"title"`
	SourceURL      string   `json:"source_url"`
	GlobalHashTags []string `json:"global_hash_tags"`
	Paragraphs     []string `json:"paragraphs"`
}
