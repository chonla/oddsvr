package api

import (
	"fmt"
	"time"
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
	activities := []Activity{}

	e := c.Get(fmt.Sprintf("%s/athlete", apiBase), &me)
	if e != nil {
		return nil, e
	}

	e = c.Get(fmt.Sprintf("%s/athletes/%d/stats", apiBase, me.ID), &stats)
	if e != nil {
		return nil, e
	}

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	endOfMonth := firstOfMonth.AddDate(0, 1, -1)
	before := endOfMonth.Unix()
	after := firstOfMonth.Unix()
	page := 1
	perPage := 100
	query := fmt.Sprintf("before=%d&after=%d&page=%d&per_page=%d", before, after, page, perPage)

	e = c.Get(fmt.Sprintf("%s/athlete/activities?%s", apiBase, query), &activities)
	if e != nil {
		return nil, e
	}

	for _, a := range activities {
		if a.Type == "Run" {
			stats.ThisMonthRunTotals.Count++
			stats.ThisMonthRunTotals.Distance += a.Distance
			stats.ThisMonthRunTotals.ElapsedTime += a.ElapsedTime
			stats.ThisMonthRunTotals.MovingTime += a.MovingTime
			stats.ThisMonthRunTotals.ElevationGain += a.ElevationGain
		}
	}

	me.Stats = stats

	return &me, nil
}
