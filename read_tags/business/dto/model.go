package dto

type Article struct {
	ID             string      `mapstructure:"id" json:"id"`
	UserID         string      `mapstructure:"user_id" json:"user_id"`
	SourceURL      string      `mapstructure:"source_url" json:"source_url"`
	GlobalHashTags []string    `mapstructure:"global_hash_tags" json:"global_hash_tags"`
	Paragraphs     []Paragraph `mapstructure:"paragraphs" json:"paragraphs"`
}

type Paragraph struct {
	ID       uint     `json:"id"`
	Text     string   `json:"text"`
	HashTags []string `json:"hash_tags"`
}

func (a *Article) ToArticleByGlobalHashTag() []ArticleByGlobalHashTag {
	var list []ArticleByGlobalHashTag
	for _, ght := range a.GlobalHashTags {
		list = append(list, ArticleByGlobalHashTag{
			GlobalHashTag: ght,
			ArticleID:     a.ID,
			SourceURL:     a.SourceURL,
		})
	}
	return list
}

func (a *Article) ToUserByArticle() UserByArticle {
	return UserByArticle{
		UserID:    a.UserID,
		ArticleID: a.ID,
	}
}

type ArticleByGlobalHashTag struct {
	GlobalHashTag string `mapstructure:"global_hash_tag" json:"global_hash_tag"`
	ArticleID     string `mapstructure:"article_id" json:"article_id"`
	SourceURL     string `mapstructure:"source_url" json:"source_url"`
}

type UserByArticle struct {
	ArticleID string `mapstructure:"article_id" json:"article_id"`
	UserID    string `mapstructure:"user_id" json:"user_id"`
	SourceURL     string `mapstructure:"source_url" json:"source_url"`
}
