package app

import (
	"os/user"
	"strconv"
)

type IdentCache struct {
	users  map[string]string
	groups map[string]string
}

func NewIdentCache() *IdentCache {
	return &IdentCache{users: map[string]string{}, groups: map[string]string{}}
}

func (c *IdentCache) User(uid string) string {
	if c == nil {
		return uid
	}
	if n, ok := c.users[uid]; ok {
		return n
	}
	id, err := strconv.Atoi(uid)
	if err != nil {
		c.users[uid] = uid
		return uid
	}
	u, err := user.LookupId(strconv.Itoa(id))
	if err != nil {
		c.users[uid] = uid
		return uid
	}
	c.users[uid] = u.Username
	return u.Username
}

func (c *IdentCache) Group(gid string) string {
	if c == nil {
		return gid
	}
	if n, ok := c.groups[gid]; ok {
		return n
	}
	g, err := user.LookupGroupId(gid)
	if err != nil {
		c.groups[gid] = gid
		return gid
	}
	c.groups[gid] = g.Name
	return g.Name
}
