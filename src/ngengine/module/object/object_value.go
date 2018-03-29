package object

type CacheData struct {
	cache map[string]interface{}
}

func (c *CacheData) Init() {
	c.cache = make(map[string]interface{})
}

// 缓存kv
func (c *CacheData) Cache(key string, value interface{}) {
	c.cache[key] = value
}

// 删除key
func (c *CacheData) DeleteCache(key string) {
	delete(c.cache, key)
}

// 是否存在key
func (c *CacheData) HasKey(key string) bool {
	_, has := c.cache[key]
	return has
}

// 获取value并返回是否存在
func (c *CacheData) TryGetValue(key string) (interface{}, bool) {
	val, has := c.cache[key]
	return val, has
}

// 获取value
func (c *CacheData) Value(key string) interface{} {
	if v, has := c.cache[key]; has {
		return v
	}
	return nil
}

// 删除所有key
func (c *CacheData) ClearAllCache() {
	c.cache = make(map[string]interface{})
}
