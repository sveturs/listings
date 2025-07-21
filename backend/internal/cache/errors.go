package cache

import "errors"

// ErrCacheMiss означает что ключ не найден в кеше
var ErrCacheMiss = errors.New("cache miss")
