package config

import "errors"

type ContextKey string

const ContextUserIDKey ContextKey = "userID"

var ErrCacheMiss = errors.New("cache: key not found")
