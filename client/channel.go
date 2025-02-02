package client

import (
	"fmt"
	"kaero/utils"
	"sort"
	"strings"
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

func (c *Channel) NicksListByMode() []string {
	nicks := make([]string, 0, len(c.Nicks))
	for k := range c.Nicks {
		nicks = append(nicks, k)
	}

	return sortWithPrefixPriority(nicks)
}

func sortWithPrefixPriority(slice []string) []string {
	priorities := map[string]int{
		"~": 1,
		"&": 2,
		"@": 3,
		"%": 4,
		"+": 5,
	}

	sort.Slice(slice, func(i, j int) bool {
		a := slice[i]
		b := slice[j]

		aPriority := 6
		bPriority := 6

		for prefix, priority := range priorities {
			if strings.HasPrefix(a, prefix) {
				aPriority = priority
			}
			if strings.HasPrefix(b, prefix) {
				bPriority = priority
			}
		}

		if aPriority != bPriority {
			return aPriority < bPriority
		} else {
			return a < b
		}
	})
	return slice
}
