package db

import (
	"blachat-server/config"
	"github.com/fatih/color"
	"github.com/gocql/gocql"
)

var session *gocql.Session

func Init() {
	var err error

	cluster := gocql.NewCluster(config.GetConfig().GetString("db_host"))

	auth := gocql.PasswordAuthenticator{
		Username: config.GetConfig().GetString("db_username"),
		Password: config.GetConfig().GetString("db_password"),
	}
	cluster.Authenticator = auth
	cluster.Keyspace = config.GetConfig().GetString("db_keyspace")
	session, err = cluster.CreateSession()

	if err != nil {
		color.Red(err.Error())
		panic(err)
	}

	color.Green("Cassandra is connected")
}

func GetSession() *gocql.Session {
	return session
}