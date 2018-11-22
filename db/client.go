package db

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Client is mongodb client
type Client struct {
	session *mgo.Session
	db      *mgo.Database
}

var dbClient *Client

// NewClient creates a new mongodb client
func NewClient(connection, database string) (*Client, error) {
	if dbClient == nil {
		session, e := mgo.Dial(connection)
		if e != nil {
			return nil, e
		}
		dbClient = &Client{
			session: session,
			db:      session.DB(database),
		}
	}
	return dbClient, nil
}

// Save data to given collection
func (c *Client) Save(collection string, data interface{}) error {
	return c.db.C(collection).Insert(data)
}

// Replace data with given one
func (c *Client) Replace(collection string, id interface{}, data interface{}) error {
	updateData := bson.M{
		"$set": data,
	}
	_, e := c.db.C(collection).UpsertId(id, updateData)
	return e
}

// Insert data with given object
func (c *Client) Insert(collection string, data interface{}) error {
	e := c.db.C(collection).Insert(data)
	return e
}
