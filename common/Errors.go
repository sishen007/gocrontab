package common

import "errors"

var (
	ERR_LOCK_ALREADY_REQUIRED = errors.New("锁已被占用")
)
