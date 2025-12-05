package domain

import (
	"strings"

	"main/internal/api"
	"main/internal/cache"
)

type ProjectService struct {
	apiClient *api.Client
	cache     *cache.Cache
}

func NewProjectService(client *api.Client, cache *cache.Cache) *ProjectService {
	return &ProjectService{
		apiClient: client,
		cache:     cache,
	}
}

func (s *ProjectService) GetAllProjects() ([]api.Project, error) {
	if projects, ok := s.cache.GetProjects(); ok {
		return projects, nil
	}

	projects, err := s.apiClient.GetProjects()
	if err != nil {
		return nil, err
	}

	s.cache.SetProjects(projects)
	return projects, nil
}

func (s *ProjectService) GetProjectByID(id string) (*api.Project, error) {
	projects, err := s.GetAllProjects()
	if err != nil {
		return nil, err
	}

	for i := range projects {
		if projects[i].ID == id {
			return &projects[i], nil
		}
	}

	return s.apiClient.GetProjectByID(id)
}

func (s *ProjectService) GetTasksForProject(projectID string) ([]api.Task, error) {
	if tasks, ok := s.cache.GetTasks(projectID); ok {
		return tasks, nil
	}

	tasks, err := s.apiClient.GetTasksForProject(projectID)
	if err != nil {
		return nil, err
	}

	s.cache.SetTasks(projectID, tasks)
	return tasks, nil
}

func (s *ProjectService) SearchProjects(query string) ([]api.Project, error) {
	projects, err := s.GetAllProjects()
	if err != nil {
		return nil, err
	}

	if query == "" {
		return projects, nil
	}

	query = strings.ToLower(query)
	var filtered []api.Project
	for _, p := range projects {
		if strings.Contains(strings.ToLower(p.Name), query) {
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}

func (s *ProjectService) RefreshCache() error {
	s.cache.Clear()
	_, err := s.GetAllProjects()
	return err
}
