package httpcache

import (
	"time"

	"github.com/chonla/oddsvr/db"
	"github.com/globalsign/mgo/bson"
)

type Cache struct {
	dbc *db.Client
}

func NewCache(dbConnection *db.Client) *Cache {
	return &Cache{
		dbc: dbConnection,
	}
}

func (c *Cache) Get(key string) (string, error) {
	m := map[string]interface{}{}

	e := c.dbc.GetBy("cache", bson.M{
		"_id": key,
	}, m)

	if e != nil {
		return "", e
	}

	return m["data"].(string), nil
}

func (c *Cache) Set(key string, maxAge string, data string) error {
	expiry := time.Now()

	d, _ := time.ParseDuration("87660h")
	if maxAge != "" {
		d, _ = time.ParseDuration(maxAge)
	}
	expiry = expiry.Add(d)

	doc := bson.M{
		"_id":    key,
		"expiry": expiry.Unix(),
		"data":   data,
	}
	return c.dbc.Insert("cache", doc)
}

func (c *Cache) Has(key string) bool {
	return c.dbc.HasBy("cache", bson.M{
		"_id": key,
		"expiry": bson.M{
			"$gte": time.Now().Unix(),
		},
	})
}
