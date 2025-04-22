package utils

import "time"

// frontend
const (
	MAX_PAGE_SIZE = 10
)

const (
	USER_NICK_NAME_PREFIX = "user_"
)

// redis
const (
	LOGIN_CODE_KEY       = "login:code:"
	LOGIN_CODE_TTL       = 2 * time.Minute
	LOGIN_TOKEN_KEY      = "login:token:"
	LOGIN_TOKEN_TTL      = 24 * time.Hour
	CACHE_SHOP_KEY       = "cache:shop:"
	CACHE_SHOP_TTL       = 30 * time.Minute
	CACHE_NULL_TTL       = 30 * time.Minute
	CACHE_SHOP_TYPE_KEY  = "cache:shop_type:"
	CACHE_SHOP_TYPE_TTL  = 0
	CACHE_SHOP_MUTEX_KEY = "lock:shop:"
	SECKILL_STOCK_KEY    = "seckill:stock:"
	LIKE_BLOG_KEY        = "blog:liked:"
	FLLOW_KEY            = "follow:"
	FEED_KEY             = "feed:"
)
