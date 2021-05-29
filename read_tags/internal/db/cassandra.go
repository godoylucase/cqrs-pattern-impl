package db

import (
	"github.com/gocql/gocql"
)

type Client struct {
	Session *gocql.Session
}

const (
	port     = 9042
	user     = "cassandra"
	password = "cassandra"
)

var (
	addresses = []string{"localhost"}
)

func Cassandra() (*Client, error) {
	clusterConfig := gocql.NewCluster(addresses...)

	clusterConfig.Authenticator = gocql.PasswordAuthenticator{Username: user, Password: password}
	clusterConfig.Port = port
	clusterConfig.ProtoVersion = 4

	session, err := clusterConfig.CreateSession()
	if err != nil {
		return nil, err
	}

	if err := configSchemaAndTables(session); err != nil {
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
