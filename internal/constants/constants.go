package constants

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// Time Pattern
	TimeFormatDateTime = "2006-01-02 15:04:05"
	TimeFormatDate     = "2006-01-02"
	TimeFormatTime     = "15:04:05"
	TimeFormatISO      = "2006-01-02T15:04:05Z"
	TimeFormatRFC3339  = "2006-01-02T15:04:05Z07:00"
	TimeFormatUnix     = "2006-01-02 15:04:05 -0700 MST"
)

var DEFAULT_JWT_SIGNING_METHOD = jwt.SigningMethodES256
var DEFAULT_JWT_EXPIRY_HOURS = time.Now().Add(time.Hour * 24)
