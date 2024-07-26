package redis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v8"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/domain/repository"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

var ErrCacheMiss = errors.New("cache: key not found")

type collectionRepository struct {
	client *redis.Client
}

func NewCollectionRepository(client *redis.Client) repository.CollectionCacheRepository {
	return &collectionRepository{
		client: client,
	}
}

func (c *collectionRepository) Get(ctx context.Context, key string) (model.Collections, error) {
	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		log.Warn("Cache miss", log.Fstring("key", key))
		return nil, ErrCacheMiss
	} else if err != nil {
		log.Error("Failed to get cache", log.Ferror(err))
		return nil, err
	}
	collections, err := c.deserialize(val)
	if err != nil {
		log.Error("Failed to deserialize collection", log.Ferror(err))
		return nil, err
	}
	log.Info("Cache hit", log.Fstring("key", key))
	return collections, nil
}

func (c *collectionRepository) Create(ctx context.Context, key string, collection model.Collections) error {
	serializeCollection, err := c.serialize(collection)
	if err != nil {
		log.Error("Failed to serialize Collection", log.Ferror(err))
		return err
	}
	if err = c.client.Set(ctx, key, serializeCollection, 0).Err(); err != nil {
		log.Error("Failed to set cache", log.Ferror(err))
		return err
	}
	log.Info("Cache set successfully", log.Fstring("key", key))
	return nil
}

func (c *collectionRepository) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		log.Error("Failed to delete cache", log.Ferror(err))
		return err
	}
	log.Info("Cache deleted successfully", log.Fstring("key", key))
	return nil
}

func (c *collectionRepository) serialize(collections model.Collections) (string, error) {
	data, err := json.Marshal(collections)
	if err != nil {
		log.Error("Failed to serialize collections", log.Ferror(err))
		return "", err
	}
	return string(data), nil
}

func (c *collectionRepository) deserialize(data string) (model.Collections, error) {
	var collections model.Collections
	err := json.Unmarshal([]byte(data), &collections)
	if err != nil {
		log.Error("Failed to deserialize collections", log.Ferror(err))
		return nil, err
	}
	return collections, nil
}
