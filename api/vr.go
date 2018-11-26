package api

import (
	"github.com/chonla/rnd"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

// NewVr creates a new fresh virtual run
func NewVr() *VirtualRun {
	return &VirtualRun{
		Engagements: []Engagement{},
	}
}

// NewVrFromContext creates a new virtual run from incoming post request
func NewVrFromContext(c echo.Context) (*VirtualRunCreateRequest, error) {
	vr := new(VirtualRunCreateRequest)
	if err := c.Bind(vr); err != nil {
		return nil, err
	}
	return vr, nil
}

// NewDistanceFromContext creates a distance from incoming post request
func NewDistanceFromContext(c echo.Context) (*Distance, error) {
	distance := new(Distance)
	if err := c.Bind(distance); err != nil {
		return nil, err
	}
	return distance, nil
}

// Save to save current object to db
func (a *API) saveVr(v *VirtualRun) (string, error) {
	var e error
	e = a.dbc.Replace("virtualrun", v.ID, v)
	return v.ID.Hex(), e
}

func (a *API) hasVr(id string) bool {
	return a.dbc.Has("virtualrun", id)
}

func (a *API) hasVrByLink(id string) bool {
	return a.dbc.HasBy("virtualrun", bson.M{
		"link": id,
	})
}

func (a *API) loadVr(id string, output *VirtualRun) error {
	return a.dbc.Get("virtualrun", id, output)
}

func (a *API) loadVrByLink(id string, output *VirtualRun) error {
	return a.dbc.GetBy("virtualrun", bson.M{
		"link": id,
	}, output)
}

func (a *API) loadMyVr(myid uint32, output *[]VirtualRun) error {
	filter := bson.M{
		"engagements": bson.M{
			"$elemMatch": bson.M{
				"athlete": myid,
			},
		},
	}

	return a.dbc.List("virtualrun", filter, output)
}

func (a *API) createSafeVrLink() string {
	link := rnd.Alphanum(12)
	return link
}

// func (a *API) loadMyVrSummary(myid uint32, output *[]VirtualRunSummary) error {
// 	aggregate := []bson.M{
// 		bson.M{
// 			"$match": bson.M{
// 				"engagements": bson.M{
// 					"$elemMatch": bson.M{
// 						"athlete": myid,
// 					},
// 				},
// 			},
// 		},
// 		bson.M{
// 			"$group": bson.M{
// 				"_id": "$_id",
// 				"count": bson.M{
// 					"$sum": 1,
// 				},
// 			},
// 		},
// 	}

// 	return a.dbc.Aggregate("virtualrun", aggregate, output)
// }
