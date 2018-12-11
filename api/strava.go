package api

import (
	"fmt"
	"time"

	"github.com/chonla/oddsvr/httpcache"
)

type Strava struct {
	c *httpcache.Cache
}

// NewStrava creates new strava interface
func NewStrava(c *httpcache.Cache) *Strava {
	return &Strava{
		c: c,
	}
}

func (s *Strava) TokenExchange(t TokenEx) (*Token, error) {
	c := client{}
	token := Token{}

	e := c.Post(apiOAuth, t, &token)
	if e != nil {
		return nil, e
	}

	return &token, nil
}

func (s *Strava) Me(token string) (*Athlete, error) {
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

	// Due to limit of response, too large response will be dropped without reason.
	// Try to request year stats by looping month-by-month instead.
	firstOfYear := time.Date(currentYear, 1, 1, 0, 0, 0, 0, currentLocation)
	firstOfYearCursor := firstOfYear
	after := firstOfYearCursor.Unix()
	var thisMonth = int(currentMonth)

	for m, lastMonth := 0, thisMonth-1; m < lastMonth; m++ {
		endOfCursor := firstOfYearCursor.AddDate(0, 1, -1)
		before := endOfCursor.Unix()

		activities, e = s.MyRunnings(&c, me.ID, before, after, 100, true)

		for _, a := range activities {
			stats.ThisYearRunTotals.Count++
			stats.ThisYearRunTotals.Distance += a.Distance
			stats.ThisYearRunTotals.ElapsedTime += a.ElapsedTime
			stats.ThisYearRunTotals.MovingTime += a.MovingTime
			stats.ThisYearRunTotals.ElevationGain += a.ElevationGain
		}

		firstOfYearCursor = firstOfYearCursor.AddDate(0, 1, 0)
		after = firstOfYearCursor.Unix()
	}

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	endOfMonth := firstOfMonth.AddDate(0, 1, -1)
	before := endOfMonth.Unix()
	after = firstOfMonth.Unix()

	activities, e = s.MyRunnings(&c, me.ID, before, after, 100, false)

	// Reuse recent activities for this month activities
	for _, a := range activities {
		stats.ThisYearRunTotals.Count++
		stats.ThisYearRunTotals.Distance += a.Distance
		stats.ThisYearRunTotals.ElapsedTime += a.ElapsedTime
		stats.ThisYearRunTotals.MovingTime += a.MovingTime
		stats.ThisYearRunTotals.ElevationGain += a.ElevationGain

		stats.ThisMonthRunTotals.Count++
		stats.ThisMonthRunTotals.Distance += a.Distance
		stats.ThisMonthRunTotals.ElapsedTime += a.ElapsedTime
		stats.ThisMonthRunTotals.MovingTime += a.MovingTime
		stats.ThisMonthRunTotals.ElevationGain += a.ElevationGain
	}

	if len(activities) > 0 {
		stats.RecentRun.Title = activities[0].Title
		stats.RecentRun.Distance = activities[0].Distance
		stats.RecentRun.MovingTime = activities[0].MovingTime
		stats.RecentRun.ElapsedTime = activities[0].ElapsedTime
		stats.RecentRun.StartDate = activities[0].StartDate
		stats.RecentRun.TimeZoneOffset = activities[0].TimeZoneOffset
	}

	me.Stats = stats

	return &me, nil
}

func (s *Strava) MyRunnings(c *client, myid uint32, before, after int64, maxResult int32, cacheFirst bool) ([]Activity, error) {
	activities := []Activity{}
	out := []Activity{}

	page := 1
	perPage := maxResult
	query := fmt.Sprintf("before=%d&after=%d&page=%d&per_page=%d", before, after, page, perPage)

	if cacheFirst {
		cacheKey := fmt.Sprintf("athlete_activities_%d_%d_%d", myid, before, after)
		c.Cacher = s.c

		e := c.GetWithCache(cacheKey, fmt.Sprintf("%s/athlete/activities?%s", apiBase, query), "", &activities)
		if e != nil {
			return nil, e
		}
	} else {
		e := c.Get(fmt.Sprintf("%s/athlete/activities?%s", apiBase, query), &activities)
		if e != nil {
			return nil, e
		}
	}

	for _, a := range activities {
		if a.Type == "Run" {
			out = append(out, a)
		}
	}

	return out, nil
}
