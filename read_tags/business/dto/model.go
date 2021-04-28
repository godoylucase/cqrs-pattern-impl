package dto

type ArticleDTO struct {
	ID             string   `mapstructure:"id" json:"id"`
	UserID         string   `mapstructure:"user_id" json:"user_id"`
	SourceURL      string   `mapstructure:"source_url" json:"source_url"`
	GlobalHashTags []string `mapstructure:"global_hash_tags" json:"global_hash_tags"`
}
