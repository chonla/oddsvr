package api

import (
	"fmt"
)

type strava struct{}

// NewStrava creates new strava interface
func NewStrava() *strava {
	return &strava{}
}

func (s *strava) TokenExchange(t TokenEx) (*Token, error) {
	c := client{}
	token := Token{}

	e := c.Post(apiOAuth, t, &token)
	if e != nil {
		return nil, e
	}

	return &token, nil
}

func (s *strava) Me(token string) (*Athlete, error) {
	c := client{
		AccessToken: token,
	}
	me := Athlete{}
	stats := Stats{}

	e := c.Get(fmt.Sprintf("%s/athlete", apiBase), &me)
	if e != nil {
		fmt.Println(e)
		return nil, e
	}

	e = c.Get(fmt.Sprintf("%s/athletes/%d/stats", apiBase, me.ID), &stats)
	if e != nil {
		fmt.Println(e)
		return nil, e
	}

	me.Stats = stats

	return &me, nil
}
