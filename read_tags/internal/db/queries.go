package db

const (
	dropKeyspaceQuery   = "DROP KEYSPACE IF EXISTS %s"

	createKeyspaceQuery = "CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy', 'datacenter1' : 1 }"

	articleByGlobalHashTagCreateQuery = `CREATE TABLE article.articles_by_global_hash_tag (
    global_hash_tag text,
    article_id text,
    source_url text,
    PRIMARY KEY ((global_hash_tag), article_id)) WITH comment = 'Q1. Finds articles by global hash tags'
	AND CLUSTERING ORDER BY (article_id DESC);`

	usersByArticleCreateQuery = `CREATE TABLE user.users_by_article (
    source_url text,
    article_id text,
    user_id text,
    PRIMARY KEY ((article_id), user_id)) WITH comment = 'Q2. Finds users by articles'
	AND CLUSTERING ORDER BY (user_id DESC);`

)
