package client

import (
	"fmt"
	"kaero/utils"
	"sync"
)

type Channel struct {
	Name  string
	Topic string

	mutex sync.Mutex
	Logs  *utils.Logger
	Nicks map[string]bool
}

func newChannel(name string) *Channel {
	return &Channel{
		Name:  name,
		Logs:  utils.NewLogger(),
		Nicks: make(map[string]bool),
	}
}

func (c *Channel) userMessage(nick string, text string) {
	c.Logs.Append(nick, utils.LogPrivMsg, text)
}

func (c *Channel) userJoin(nick string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Nicks[nick] = true
	// @todo: handle prefixes
	c.Logs.Append(nick, utils.LogJoined, "joined.")
}

func (c *Channel) usersJoin(nicks []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, nick := range nicks {
		c.Nicks[nick] = true
		// @todo: handle prefixes
	}
}

func (c *Channel) userLeave(nick string, reason string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.Nicks[nick]; ok {
		delete(c.Nicks, nick)
		text := "left."
		if reason != "" {
			text = fmt.Sprintf("left. <%s>", reason)
		}
		c.Logs.Append(nick, utils.LogLeft, text)
	}
}

func (c *Channel) userPart(nick string, reason string) {
	c.userLeave(nick, reason)
}

func (c *Channel) userQuit(nick string, reason string) {
	c.userLeave(nick, reason)
}

func (c *Channel) userNick(oldNick string, newNick string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.Nicks[oldNick]; ok {
		// @todo: handle prefixes
		delete(c.Nicks, oldNick)
		c.Nicks[newNick] = true
		text := fmt.Sprintf("%s changed their nick to %s.", oldNick, newNick)
		c.Logs.Append(oldNick, utils.LogSystem, text)
	}
}
