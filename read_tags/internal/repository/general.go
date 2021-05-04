package repository

import (
	"fmt"
	"strings"

	"github.com/godoylucase/read_tags/business/dto"
	"github.com/godoylucase/read_tags/internal/db"
	"github.com/mitchellh/mapstructure"
)

type repo struct {
	client *db.Client
}

func NewArticle(client *db.Client) *repo {
	return &repo{client: client}
}

func (r *repo) SaveArticleByGlobalTags(dto dto.ArticleByGlobalHashTag) error {
	query := fmt.Sprintf(
		"INSERT INTO %s.%s (global_hash_tag, article_id , source_url) VALUES (?,?,?)",
		db.ArticleSpace,
		db.ArticleByGlobalHashTagsTable)

	if err := r.client.Session.Query(query, dto.GlobalHashTag, dto.ArticleID, dto.SourceURL).Exec(); err != nil {
		return err
	}

	return nil
}

func (r *repo) SaveUserByArticle(dto dto.UserByArticle) error {
	query := fmt.Sprintf(
		"INSERT INTO %s.%s (article_id, user_id , source_url) VALUES (?,?,?)",
		db.UserSpace,
		db.UserByArticleTable)

	if err := r.client.Session.Query(query, dto.ArticleID, dto.UserID, dto.SourceURL).Exec(); err != nil {
		return err
	}

	return nil
}

func (r *repo) GetArticleByGlobalTags(globalHashTags []string) ([]dto.ArticleByGlobalHashTag, error) {
	var ght []string
	for _, value := range globalHashTags {
		v := fmt.Sprintf("'%s'", value)
		ght = append(ght, v)
	}

	in := strings.Join(ght, ",")
	query := fmt.Sprintf("SELECT * FROM %s.%s where global_hash_tag IN (%s)", db.ArticleSpace, db.ArticleByGlobalHashTagsTable, in)

	results, err := r.client.Session.Query(query).Iter().SliceMap()
	if err != nil {
		return nil, err
	}

	var articles []dto.ArticleByGlobalHashTag
	var article dto.ArticleByGlobalHashTag
	for _, r := range results {
		if err := mapstructure.Decode(r, &article); err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (r *repo) GetUsersByArticle(articleID string) ([]dto.UserByArticle, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s where article_id=%s", db.ArticleSpace, db.ArticleByGlobalHashTagsTable, articleID)

	results, err := r.client.Session.Query(query).Iter().SliceMap()
	if err != nil {
		return nil, err
	}

	var usersByArticle []dto.UserByArticle
	var userByArticle dto.UserByArticle
	for _, r := range results {
		if err := mapstructure.Decode(r, &userByArticle); err != nil {
			return nil, err
		}

		usersByArticle = append(usersByArticle, userByArticle)
	}

	return usersByArticle, nil
}
