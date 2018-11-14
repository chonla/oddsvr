package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Client is mongodb client
type Client struct {
	session *mgo.Session
	db      *mgo.Database
}

// NewClient creates a new mongodb client
func NewClient(connection, database string) (*Client, error) {
	session, e := mgo.Dial(connection)
	if e != nil {
		return nil, e
	}
	return &Client{
		session: session,
		db:      session.DB(database),
	}, nil
}

// Save data to given collection
func (c *Client) Save(collection string, data interface{}) error {
	return c.db.C(collection).Insert(data)
}

// Replace data with given one
func (c *Client) Replace(collection string, id uint32, data interface{}) error {
	updateData := bson.M{
		"$set": data,
	}
	_, e := c.db.C(collection).UpsertId(id, updateData)
	return e
}
