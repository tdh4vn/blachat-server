package repo

import (
	"github.com/go-redis/redis"
)

type PresenceRepo interface {
	SetUserOnline(uuid string) (int64, error)
	SetUserOffline(uuid string) (int64, error)
	CheckUsersOnline(uuid []string) []string
}

func NewPresenceRepo(redisClient *redis.Client) PresenceRepo{
	return &PresenceRepoImpl{
		redisClient: redisClient,
	}
}

type PresenceRepoImpl struct {
	redisClient *redis.Client
}

func (presenceRepo *PresenceRepoImpl) SetUserOnline(uuid string) (int64, error) {
	return presenceRepo.redisClient.Incr(uuid).Result()
}

func (presenceRepo *PresenceRepoImpl) SetUserOffline(uuid string) (int64, error) {
	return presenceRepo.redisClient.Decr(uuid).Result()
}

func (presenceRepo *PresenceRepoImpl) CheckUsersOnline(uuid []string) []string {
	var results []string

	for _, id := range uuid {
		if numberClient, err := presenceRepo.redisClient.Get(id).Int64(); err == nil && numberClient > 0 {
			results = append(results, id)
		}
	}

	return results
}


