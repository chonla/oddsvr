package api

import (
	"github.com/globalsign/mgo/bson"
)

// Athlete is runner profile
type Athlete struct {
	ID                   uint32 `json:"id"`
	UserName             string `json:"username"`
	FirstName            string `json:"firstname"`
	LastName             string `json:"lastname"`
	City                 string `json:"city"`
	State                string `json:"state"`
	Country              string `json:"country"`
	Gender               string `json:"sex"`
	ProfilePicture       string `json:"profile"`
	ProfilePictureMedium string `json:"profile_medium"`
	Email                string `json:"email"`
	Stats                `json:"stats"`
}

// Stats is running stats
type Stats struct {
	RecentRunTotals    RunStats `json:"recent_run_totals"`
	AllRunTotals       RunStats `json:"all_run_totals"`
	ThisMonthRunTotals RunStats `json:"this_month_run_totals"`
}

// RunStats is detailed of stats
type RunStats struct {
	Count         uint32  `json:"count"`
	Distance      float64 `json:"distance"`
	MovingTime    uint32  `json:"moving_time"`
	ElapsedTime   uint32  `json:"elapsed_time"`
	ElevationGain float64 `json:"elevation_gain"`
}

// Activity is activity
type Activity struct {
	Distance       float64 `json:"distance"`
	MovingTime     uint32  `json:"moving_time"`
	ElapsedTime    uint32  `json:"elapsed_time"`
	ElevationGain  float64 `json:"total_elevation_gain"`
	Type           string  `json:"type"`
	StartDate      string  `json:"string"`
	TimeZoneOffset float64 `json:"utc_offset"`
}

// VirtualRun is virtual run
type VirtualRun struct {
	ID       bson.ObjectId `json:"_id,omitempty"`
	Name     string        `json:"name"`
	FromDate string        `json:"from_date"`
	ToDate   string        `json:"to_date"`
}
