package api

func tokenRefresh(t TokenRefresh) (*Token, error) {
	c := client{}
	token := Token{}

	e := c.Post(apiOAuth, t, &token)
	if e != nil {
		return nil, e
	}

	return &token, nil
}
