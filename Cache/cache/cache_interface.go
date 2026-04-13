package cache

type Cache interface {
	Get(key string) (any, bool)
	Put(key string, value any)
	Delete(key string) bool
	Size() int
	Capacity() int
	Clear()
	String() string
}
