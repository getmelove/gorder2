package factory

import "sync"

// Supplier 定义了一个函数类型，用于在缓存未命中时创建新的实例。
// 它接收一个字符串 key 并返回任意类型的值 (any)。
type Supplier func(string) any

// Singleton 结构体实现了一个通用的单例工厂模式。
// 它维护一个 map 来缓存创建的对象，并使用 mutex 保证线程安全。
type Singleton struct {
	// cache 用于存储已创建的对象实例，key 为标识符，value 为对象实例。
	cache map[string]any
	// locker 是互斥锁，用于保护 cache map 在并发访问时的安全性。
	locker *sync.Mutex
	// supplier 是当缓存中找不到对象时，用于创建新对象的函数。
	supplier Supplier
}

// NewSingleton 创建并返回一个新的 Singleton 实例。
// 需要传入一个 Supplier 函数，用于后续创建具体的对象。
func NewSingleton(supplier Supplier) *Singleton {
	return &Singleton{
		cache:    make(map[string]any),
		locker:   &sync.Mutex{},
		supplier: supplier,
	}
}

// Get 根据 key 获取对象实例。
// 这是一个线程安全的方法，采用"双重检查锁定"（Double-Checked Locking）的变体思想。
func (s *Singleton) Get(key string) any {
	// 第一次检查：如果在不加锁的情况下已经能从缓存读到，直接返回，避免锁的开销。
	if value, hit := s.cache[key]; hit {
		return value
	}

	// 加锁：确保只有一个 goroutine 能执行后续的创建逻辑。
	s.locker.Lock()
	defer s.locker.Unlock()

	// 第二次检查：获取锁之后再次检查缓存。
	// 这是为了防止在获取锁的过程中，其他 goroutine 已经创建了该对象。
	if value, hit := s.cache[key]; hit {
		return value
	}

	// 如果缓存中确实没有，调用 supplier 创建新对象。
	v := s.supplier(key)
	// 将新对象存入缓存。
	s.cache[key] = v

	return v
}
