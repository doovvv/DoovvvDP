package utils

import "time"
const (
	USER_NICK_NAME_PREFIX = "user_"
)

//redis
const(
	LOGIN_CODE_KEY = "login:code:"
	LOGIN_CODE_TTL = 2*time.Minute
	LOGIN_TOKEN_KEY = "login:token:"
	LOGIN_TOKEN_TTL = 24*time.Hour
)