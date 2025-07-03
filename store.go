package captcha

import (
	"runtime"
	"sync"
	"time"
)

/*MapStore MapStore */
type MapStore struct {
	store      sync.Map
	stopCh     chan struct{}
	expiration time.Duration
}

type stored struct {
	val any
	ttl time.Time
}

// NewMapStore ...
func NewMapStore(expiration, cleanupInterval time.Duration) *MapStore {
	m := &MapStore{
		expiration: expiration,
		stopCh:     make(chan struct{}),
	}
	// 只有当清理间隔大于0时，才启动后台清理 goroutine
	if cleanupInterval > 0 {
		go m.cleanupLoop(cleanupInterval)
		// 设置Finalizer，当MapStore对象被GC回收时，自动关闭stopCh，
		// 从而优雅地停止后台goroutine，防止内存泄漏。
		runtime.SetFinalizer(m, (*MapStore).stopCleanup)
	}

	return m
}

// cleanupLoop 是后台运行的清理循环
func (m *MapStore) cleanupLoop(interval time.Duration) {
	// 创建一个定时器
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		// 等待定时器触发
		case <-ticker.C:
			// 遍历 map 中的所有键值对
			m.store.Range(func(key, value any) bool {
				if s, ok := value.(*stored); ok {
					// 检查是否过期
					if !s.ttl.IsZero() && s.ttl.Before(time.Now()) {
						// 已过期，删除
						m.store.Delete(key)
					}
				}
				// 返回 true 继续遍历
				return true
			})
		// 等待停止信号
		case <-m.stopCh:
			// 收到停止信号，退出循环
			return
		}
	}
}

// stopCleanup 用于被 runtime.SetFinalizer 调用
func (m *MapStore) stopCleanup() {
	// 使用 select + default 来避免在已经关闭的 channel 上再次 close 导致 panic
	select {
	case <-m.stopCh:
		return
	default:
		close(m.stopCh)
	}
}

/*GetD get interface with default */
func (m *MapStore) Get(key string) any {
	if v, b := m.store.Load(key); b {
		m.store.Delete(key)
		switch vv := v.(type) {
		case *stored:
			if !vv.ttl.IsZero() && vv.ttl.Before(time.Now()) {
				return nil
			}
			return vv.val
		}
	}
	return nil
}

/*Set set interface with ttl */
func (m *MapStore) Set(key string, val any) {
	m.store.Store(key, &stored{
		val: val,
		ttl: time.Now().Add(m.expiration),
	})
}

/*Has check exist */
func (m *MapStore) Has(key string) bool {
	if v, b := m.store.Load(key); b {
		switch vv := v.(type) {
		case *stored:
			if !vv.ttl.IsZero() && vv.ttl.Before(time.Now()) {
				return false
			}
			return true
		}
	}
	return false
}

/*Delete one value */
func (m *MapStore) Delete(key string) {
	m.store.Delete(key)
}

/*Clear delete all values */
func (m *MapStore) Clear() {
	*m = MapStore{}
}
