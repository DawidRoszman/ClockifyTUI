package domain

import (
	"main/internal/api"
	"main/internal/cache"
)

type TagService struct {
	apiClient *api.Client
	cache     *cache.Cache
}

func NewTagService(client *api.Client, cache *cache.Cache) *TagService {
	return &TagService{
		apiClient: client,
		cache:     cache,
	}
}

func (s *TagService) GetAllTags() ([]api.Tag, error) {
	if tags, ok := s.cache.GetTags(); ok {
		return tags, nil
	}

	tags, err := s.apiClient.GetTags()
	if err != nil {
		return nil, err
	}

	s.cache.SetTags(tags)
	return tags, nil
}
