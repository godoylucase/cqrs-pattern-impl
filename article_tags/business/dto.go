package business

type ArticleDTO struct {
	ID             string      `json:"id"`
	UserID         UserID      `json:"user_id"`
	Title          string      `json:"title"`
	SourceURL      URL         `json:"source_url"`
	GlobalHashTags []HashTag   `json:"global_hash_tags"`
	Paragraphs     []Paragraph `json:"paragraphs"`
}

func (dto *ArticleDTO) GetKey() string {
	return dto.ID
}
