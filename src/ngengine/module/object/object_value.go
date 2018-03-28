package object

type CacheData struct {
	cache map[string]interface{}
}

func (c *CacheData) Init() {
	c.cache = make(map[string]interface{})
}

func (c *CacheData) Cache(key string, value interface{}) {
	c.cache[key] = value
}

func (c *CacheData) Value(key string) interface{} {
	if v, has := c.cache[key]; has {
		return v
	}
	return nil
}

func (c *CacheData) ClearAllCache() {
	c.cache = make(map[string]interface{})
}
