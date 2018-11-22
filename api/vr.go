package api

import (
	"github.com/labstack/echo"
)

// NewVr creates a new fresh virtual run
func NewVr() *VirtualRun {
	return &VirtualRun{
		Athletes: []uint32{},
	}
}

// NewVrFromContext creates a new virtual run from incoming post request
func NewVrFromContext(c echo.Context) (*VirtualRun, error) {
	vr := new(VirtualRun)
	if err := c.Bind(vr); err != nil {
		return nil, err
	}
	return vr, nil
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

func (a *API) loadVr(id string, output *VirtualRun) error {
	return a.dbc.Get("virtualrun", id, output)
}
