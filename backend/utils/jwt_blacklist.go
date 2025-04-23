package utils

import (
	"sync"
	"time"
)

var (
	tokenBlacklist = make(map[string]time.Time)
	blacklistMutex = &sync.RWMutex{}
)

// Thêm token vào danh sách đen
func AddToBlacklist(token string, expiry time.Duration) error {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	tokenBlacklist[token] = time.Now().Add(expiry)
	return nil
}

// Kiểm tra token có bị blacklist không
func IsTokenBlacklisted(token string) bool {
	blacklistMutex.RLock()
	expiry, exists := tokenBlacklist[token]
	blacklistMutex.RUnlock()

	if !exists {
		return false
	}

	// Nếu hết hạn thì xóa (dùng lock ghi)
	if time.Now().After(expiry) {
		blacklistMutex.Lock()
		delete(tokenBlacklist, token)
		blacklistMutex.Unlock()
		return false
	}

	return true
}

// GC chạy ngầm mỗi N phút để xóa token hết hạn
func StartTokenBlacklistGC(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)

			now := time.Now()
			blacklistMutex.Lock()
			for token, expiry := range tokenBlacklist {
				if now.After(expiry) {
					delete(tokenBlacklist, token)
				}
			}
			blacklistMutex.Unlock()
		}
	}()
}
