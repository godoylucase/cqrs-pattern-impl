package repository

import (
	"fmt"
	"strings"

	"github.com/godoylucase/read_tags/business/dto"
	"github.com/godoylucase/read_tags/internal/db"
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

func (r *repo) UpsertUserBySourceURL(dto dto.Article) error {
	query := fmt.Sprintf(
		"INSERT INTO %s.%s (source_url, user_id, article_id) VALUES (?,?,?)",
		db.UserSpace,
		db.UserArticlesBySourceURLTable)

	if err := r.client.Session.Query(query, dto.SourceURL, dto.UserID, dto.ID).Exec(); err != nil {
		return err
	}

	return nil
}

func (r *repo) ArticleByGlobalTags(globalHashTags []string) (dto.ArticleByGlobalHashTagRead, error) {
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

func (r *repo) UserArticlesBySourceURL(sourceUrl string) (dto.UserArticlesBySourceURLRead, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s where source_url='%s'", db.UserSpace, db.UserArticlesBySourceURLTable, sourceUrl)

	queryResults, err := r.client.Session.Query(query).Iter().SliceMap()
	if err != nil {
		return dto.UserArticlesBySourceURLRead{}, err
	}

	sus := make(map[string][]dto.UserArticle)
	for _, qr := range queryResults {
		key := qr["source_url"].(string)

		read := dto.UserArticle{
			ArticleID: qr["article_id"].(string),
			UserID:    qr["user_id"].(string),
		}

		existingDets, ok := sus[key]
		if !ok {
			var ua []dto.UserArticle
			sus[key] = append(ua, read)
			continue
		}

		existingDets = append(existingDets, read)
		sus[key] = existingDets
	}

	return dto.UserArticlesBySourceURLRead{SourceURLs: sus}, nil
}
