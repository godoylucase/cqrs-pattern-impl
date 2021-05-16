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

func (r *repo) UpsertArticleByGlobalTags(dto dto.ArticleByGlobalHashTag) error {
	query := fmt.Sprintf(
		"INSERT INTO %s.%s (global_hash_tag, article_id , source_url) VALUES (?,?,?)",
		db.ArticleSpace,
		db.ArticleByGlobalHashTagsTable)

	if err := r.client.Session.Query(query, dto.GlobalHashTags, dto.Detail.ArticleID, dto.Detail.SourceURL).Exec(); err != nil {
		return err
	}

	return nil
}

func (r *repo) UpsertUserByArticle(dto dto.UserByArticle) error {
	query := fmt.Sprintf(
		"INSERT INTO %s.%s (article_id, user_id , source_url) VALUES (?,?,?)",
		db.UserSpace,
		db.UserByArticleTable)

	if err := r.client.Session.Query(query, dto.ArticleID, dto.UserID, dto.SourceURL).Exec(); err != nil {
		return err
	}

	return nil
}

func (r *repo) GetArticleByGlobalTags(globalHashTags []string) (dto.ArticleByGlobalHashTagRead, error) {
	var ght []string
	for _, value := range globalHashTags {
		v := fmt.Sprintf("'%s'", value)
		ght = append(ght, v)
	}

	in := strings.Join(ght, ",")
	query := fmt.Sprintf("SELECT * FROM %s.%s where global_hash_tag IN (%s)", db.ArticleSpace, db.ArticleByGlobalHashTagsTable, in)

	queryResults, err := r.client.Session.Query(query).Iter().SliceMap()
	if err != nil {
		return dto.ArticleByGlobalHashTagRead{}, err
	}

	resultMap := make(map[string][]dto.ArticleByGlobalHashTagDetail)
	for _, qr := range queryResults {
		key := qr["global_hash_tag"].(string)

		hashTagDetail := dto.ArticleByGlobalHashTagDetail{
			ArticleID: qr["article_id"].(string),
			SourceURL: qr["source_url"].(string),
		}

		existingDets, ok := resultMap[key]
		if !ok {
			var details []dto.ArticleByGlobalHashTagDetail
			resultMap[key] = append(details, hashTagDetail)
			continue
		}

		existingDets = append(existingDets, hashTagDetail)
		resultMap[key] = existingDets
	}

	return dto.ArticleByGlobalHashTagRead{GlobalHashTags: resultMap}, nil
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
