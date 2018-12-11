package api

import (
	"fmt"
	"time"

	"github.com/kr/pretty"
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

	// Due to limit of response, too large response will be dropped without reason.
	// Try to request year stats by looping month-by-month instead.
	firstOfYear := time.Date(currentYear, 1, 1, 0, 0, 0, 0, currentLocation)
	firstOfYearCursor := firstOfYear
	after := firstOfYearCursor.Unix()
	var thisMonth = int(currentMonth)
	for m := 0; m < thisMonth; m++ {
		endOfCursor := firstOfYearCursor.AddDate(0, 1, -1)
		before := endOfCursor.Unix()

		activities, e = s.MyRunnings(&c, before, after, 100)

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

	// Reuse recent activities for this month activities
	for _, a := range activities {
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

func (s *strava) MyRunnings(c *client, before, after int64, maxResult int32) ([]Activity, error) {
	activities := []Activity{}
	out := []Activity{}

	page := 1
	perPage := maxResult
	query := fmt.Sprintf("before=%d&after=%d&page=%d&per_page=%d", before, after, page, perPage)

	e := c.Get(fmt.Sprintf("%s/athlete/activities?%s", apiBase, query), &activities)
	if e != nil {
		return nil, e
	}

	pretty.Println(activities)

	for _, a := range activities {
		if a.Type == "Run" {
			out = append(out, a)
		}
	}

	return out, nil
}
