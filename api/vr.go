package api

import (
	"github.com/labstack/echo"
)

// NewVr creates a new fresh virtual run
func NewVr() *VirtualRun {
	return &VirtualRun{}
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
	return v.ID.String(), e
}
