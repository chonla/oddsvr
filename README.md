# ODDS Virtual Run API

## Starting by environment variable

```
<environment variables initialization...> oddsvr
```

## Supported environment

| Name | Description |
| --- | --- |
| ODDSVR_ADDR | (default 0.0.0.0:1323) Listening address |
| ODDSVR_DB | (default 127.0.0.1:27107) MongoDB server address |
| ODDSVR_STRAVA_CLIENT_ID | Client ID registered with Strava |
| ODDSVR_STRAVA_CLIENT_SECRET | Client Secret from Strava |
| ODDSVR_JWT_SECRET | JWT Secret for communicate with browser |