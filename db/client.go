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

// Has returns true if collection contain document with given id
func (c *Client) Has(collection string, id interface{}) bool {
	q := c.db.C(collection).FindId(bson.ObjectIdHex(id.(string))).Limit(1)
	count, e := q.Count()
	if e != nil {
		return false
	}
	return (count > 0)
}

// Get load data from collection with given id into output
func (c *Client) Get(collection string, id interface{}, output interface{}) error {
	q := c.db.C(collection).FindId(bson.ObjectIdHex(id.(string))).Limit(1)
	return q.Select(bson.M{}).One(output)
}

// List all data from given filter in collection into output
func (c *Client) List(collection string, filter interface{}, output interface{}) error {
	return c.db.C(collection).Find(filter).All(output)
}
