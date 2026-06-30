package cache

import (
	"context"
	"delivery-tracker/internal/domain"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type ParcelCache struct {
	client *redis.Client
}

func NewParcelCache(client *redis.Client) *ParcelCache {
	return &ParcelCache{client: client}
}

func parcelKey(trackNumber string) string {
	return fmt.Sprintf("parcel:track:%s", trackNumber)
}

func (c *ParcelCache) SetByTrack(
	ctx context.Context,
	trackNumber string,
	parcel *domain.Parcel,
	ttl time.Duration) error {

	data, err := json.Marshal(parcel)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	err = c.client.Set(ctx, parcelKey(trackNumber), data, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}

func (c *ParcelCache) GetByTrack(ctx context.Context, trackNumber string) (*domain.Parcel, error) {
	var parcel domain.Parcel
	val, err := c.client.Get(ctx, parcelKey(trackNumber)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("cache get parcel by track: %w", err)
	}

	if err = json.Unmarshal([]byte(val), &parcel); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return &parcel, nil
}

func (c *ParcelCache) DeleteByTrack(ctx context.Context, trackNumber string) error {
	err := c.client.Del(ctx, parcelKey(trackNumber)).Err()
	if err != nil {
		return fmt.Errorf("cache delete parcel by track: %w", err)
	}

	return nil
}
