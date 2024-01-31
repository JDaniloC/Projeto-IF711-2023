package utils

import (
	"fmt"
	"sync"
)

type Controller struct {
	ValidLinks   StringSet
	InvalidLinks StringSet
	VisitedLinks StringSet
	MaxDepth     int
	mu           sync.RWMutex
}

func (c *Controller) AddValidLink(link string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ValidLinks.Add(link)
}

func (c *Controller) AddInvalidLink(link string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.InvalidLinks.Add(link)
}

func (c *Controller) AddVisitedLink(link string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.VisitedLinks.Add(link)
}

func (c *Controller) CanVisit(link string, depth int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if depth >= c.MaxDepth {
		return false
	}
	if c.VisitedLinks.Contains(link) {
		return false
	}
	return true
}

func (c *Controller) String() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return fmt.Sprintf("%d x %d", c.ValidLinks.Len(), c.InvalidLinks.Len())
}

func NewController(maxDepth int) *Controller {
	controller := Controller{
		ValidLinks:   *NewStringSet(),
		InvalidLinks: *NewStringSet(),
		VisitedLinks: *NewStringSet(),
		MaxDepth:     maxDepth,
	}
	return &controller
}
