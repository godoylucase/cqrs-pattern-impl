package repository

import (
	"fmt"

	"github.com/godoylucase/read_tags/business/dto"
	"github.com/godoylucase/read_tags/internal/db"
)

type repo struct {
	client *db.Client
}

func NewArticle(client *db.Client) *repo {
	return &repo{client: client}
}

func (r repo) SaveArticleByGlobalTags(dto dto.ArticleByGlobalHashTag) error {
	query := fmt.Sprintf(
		"INSERT INTO %s.%s (global_hash_tag, article_id , source_url) VALUES (?,?,?)",
		db.ArticleSpace,
		db.ArticleByGlobalHashTagsTable)

	if err := r.client.Session.Query(query).Exec(); err != nil {
		return err
	}

	return nil
}

func (r repo) SaveUserByArticle(dto dto.UserByArticle) error {
	query := fmt.Sprintf(
		"INSERT INTO %s.%s (article_id, user_id , source_url) VALUES (?,?,?)",
		db.UserSpace,
		db.UserByArticleTable)

	if err := r.client.Session.Query(query).Exec(); err != nil {
		return err
	}

	return nil
}
