package cache

import (
	"sync"
	"time"

	"main/internal/api"
)

type Cache struct {
	projects      []api.Project
	tasks         map[string][]api.Task
	projectsMutex sync.RWMutex
	tasksMutex    sync.RWMutex
	ttl           time.Duration
	lastUpdate    time.Time
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		tasks:      make(map[string][]api.Task),
		ttl:        ttl,
		lastUpdate: time.Time{},
	}
}

func (c *Cache) SetProjects(projects []api.Project) {
	c.projectsMutex.Lock()
	defer c.projectsMutex.Unlock()
	c.projects = projects
	c.lastUpdate = time.Now()
}

func (c *Cache) GetProjects() ([]api.Project, bool) {
	c.projectsMutex.RLock()
	defer c.projectsMutex.RUnlock()

	if c.IsExpired() || len(c.projects) == 0 {
		return nil, false
	}

	return c.projects, true
}

func (c *Cache) SetTasks(projectID string, tasks []api.Task) {
	c.tasksMutex.Lock()
	defer c.tasksMutex.Unlock()
	c.tasks[projectID] = tasks
}

func (c *Cache) GetTasks(projectID string) ([]api.Task, bool) {
	c.tasksMutex.RLock()
	defer c.tasksMutex.RUnlock()

	if c.IsExpired() {
		return nil, false
	}

	tasks, ok := c.tasks[projectID]
	return tasks, ok
}

func (c *Cache) IsExpired() bool {
	if c.lastUpdate.IsZero() {
		return true
	}
	return time.Since(c.lastUpdate) > c.ttl
}

func (c *Cache) Clear() {
	c.projectsMutex.Lock()
	c.tasksMutex.Lock()
	defer c.projectsMutex.Unlock()
	defer c.tasksMutex.Unlock()

	c.projects = nil
	c.tasks = make(map[string][]api.Task)
	c.lastUpdate = time.Time{}
}
