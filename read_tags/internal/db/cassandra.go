package db

import (
	"crypto/tls"
	"strconv"

	"github.com/gocql/gocql"
)

type Client struct {
	Session *gocql.Session
}

const (
	contactPoint = ""
	port         = ""
	user         = ""
	password     = ""
)

func Cassandra() (*Client, error) {
	clusterConfig := gocql.NewCluster(contactPoint)
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	clusterConfig.Authenticator = gocql.PasswordAuthenticator{Username: user, Password: password}
	clusterConfig.Port = p
	clusterConfig.SslOpts = &gocql.SslOptions{Config: &tls.Config{MinVersion: tls.VersionTLS12}}
	clusterConfig.ProtoVersion = 4

	session, err := clusterConfig.CreateSession()
	if err != nil {
		return nil, err
	}

	return &Client{Session: session}, nil
}

func (c *Client) Close() {
	c.Session.Close()
}

func configSchemaAndTables(session *gocql.Session) error {
	if err := initialize(session); err != nil {
		return err
	}

	return nil
}
