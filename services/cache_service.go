package services

import (
	"fmt"
	"sync"
)

type CacheService struct {
	characters     map[string][]byte
	characterMutex sync.RWMutex
	marvelService  *MarvelService
}

func NewCacheService() *CacheService {
	return &CacheService{
		characters:    make(map[string][]byte),
		marvelService: NewMarvelService(),
	}
}

func (c *CacheService) GetCharacters(offset int, limit int) ([]byte, error) {
	c.characterMutex.RLock()
	defer c.characterMutex.RUnlock()

	key := fmt.Sprintf("offset_%d", offset)
	if data, ok := c.characters[key]; ok {
		return data, nil
	}

	data, err := c.marvelService.FetchMarvelCharacters(offset, limit)
	if err != nil {
		return nil, err
	}
	return []byte(data), nil
}
