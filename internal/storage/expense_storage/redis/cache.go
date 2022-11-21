package redis

import (
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client *redis.Client
}

//func (c *Cache) GetByRange(ctx context.Context, userID int64, timeRange int) ([]*models.TotalExpense, error) {
//	c.client.Get(ctx, userID, timeRange)
//}
//
//func (c *Cache) SetByRange(ctx context.Context, userID int64, timeRange int, exps []*models.TotalExpense) error {
//
//}
