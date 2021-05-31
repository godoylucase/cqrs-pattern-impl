package db

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
)

const (
	ArticleSpace                 = "article"
	ArticleByGlobalHashTagsTable = "articles_by_global_hash_tag"
	UserSpace                    = "user"
	UserArticlesBySourceURLTable = "users_articles_by_source_url"
)

func initialize(session *gocql.Session) error {
	if err := v0_0_1__20210428_init(session); err != nil {
		return err
	}

	return nil
}

func dropKeySpaceIfExists(keyspace string, Session *gocql.Session) error {
	err := Session.Query(fmt.Sprintf(dropKeyspaceQuery, keyspace)).Exec()
	if err != nil {
		return err
	}

	logrus.Infof("keyspace %v dropped", keyspace)
	return nil
}

func createKeySpace(keyspace string, session *gocql.Session) error {
	err := session.Query(fmt.Sprintf(createKeyspaceQuery, keyspace)).Exec()
	if err != nil {
		return err
	}

	logrus.Infof("keyspace %v created", keyspace)
	return nil
}

func v0_0_1__20210428_init(session *gocql.Session) error {
	if err := dropKeySpaceIfExists(ArticleSpace, session); err != nil {
		return err
	}

	if err := dropKeySpaceIfExists(UserSpace, session); err != nil {
		return err
	}

	if err := createKeySpace(ArticleSpace, session); err != nil {
		return err
	}

	if err := createKeySpace(UserSpace, session); err != nil {
		return err
	}

	err := session.Query(articleByGlobalHashTagCreateQuery).Exec()
	if err != nil {
		return err
	}
	logrus.Infof("table %s.%s created", ArticleSpace, ArticleByGlobalHashTagsTable)

	if err := session.Query(UserArticlesBySourceURLCreateQuery).Exec(); err != nil {
		return err
	}
	logrus.Infof("table %s.%s created", UserSpace, UserArticlesBySourceURLTable)

	return nil
}
