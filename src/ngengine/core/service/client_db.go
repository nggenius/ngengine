package service

import (
	"net"
	"ngengine/share"
	"sync"
)

type ClientDB struct {
	sync.RWMutex
	ctx     *context
	clients map[uint64]*Client
	serial  uint64
	quit    bool
}

func NewClientDB(ctx *context) *ClientDB {
	db := &ClientDB{
		clients: make(map[uint64]*Client, 2048),
		serial:  0,
		ctx:     ctx,
		quit:    false,
	}
	return db
}

func (c *ClientDB) AddClient(conn net.Conn) uint64 {
	c.Lock()
	defer c.Unlock()

	if c.quit {
		return 0
	}

	if conn == nil {
		return 0
	}

	//寻找一个唯一ID
	for {
		c.serial++
		if c.serial > share.SESSION_MAX {
			c.serial = 1
		}
		if _, dup := c.clients[c.serial]; dup {
			continue
		}
		break
	}

	client := NewClient(c.serial, conn, c.ctx)
	c.clients[client.Session] = client
	c.ctx.Core.LogInfo("add ", client)
	return client.Session
}

func (c *ClientDB) FindClient(session uint64) *Client {
	c.RLock()
	defer c.RUnlock()
	if client, has := c.clients[session]; has {
		return client
	}

	return nil
}

func (c *ClientDB) BreakClient(session uint64) {
	c.RLock()
	defer c.RUnlock()
	if client, has := c.clients[session]; has {
		client.Close()
	}
}

func (c *ClientDB) RemoveClient(session uint64) {
	c.Lock()
	defer c.Unlock()
	if client, has := c.clients[session]; has {
		delete(c.clients, session)
		client.Close()
		c.ctx.Core.LogInfo("remove ", client)
	}
}

func (c *ClientDB) CloseAll() {
	c.Lock()
	defer c.Unlock()
	if !c.quit {
		c.quit = true
		for _, client := range c.clients {
			client.Close()
		}
	}

}
