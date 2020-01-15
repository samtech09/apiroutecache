package apiroutecache

import (
	"strconv"
	"time"

	mgo "gopkg.in/mgo.v2"
)

//MongoConfig is config for mongodb connection
type MongoConfig struct {
	Host        string
	Port        int
	DbName      string
	User        string
	Pwd         string
	ConnTimeout int // in seconds
}

//MongoSession holds session for mongodb connection
type MongoSession struct {
	*mgo.Session
	DBname string
}

//InitSession initialize mongodb session
func InitSession(config MongoConfig) *MongoSession {
	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{config.Host + ":" + strconv.Itoa(config.Port)},
		Timeout:  time.Duration(config.ConnTimeout) * time.Second,
		Database: config.DbName,
		Username: config.User,
		Password: config.Pwd,
	}

	// Connect to mongoDB
	var err error
	s, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic(err)
	}
	return &MongoSession{s, config.DbName}
}

//Cleanup closes existing mongodb session
func (s *MongoSession) Cleanup() {
	s.Close()
}
