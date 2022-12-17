package alock

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/go-redis/redis"
)

var RedisDB *redis.Client

func InitRedisLock(host, password string, db, poolSize int) error {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	})

	_, err := RedisDB.Ping().Result()
	return err
}

type RLock struct {
	key string
	reqID string
	duration time.Duration
}

func NewRLock(key, reqID string, duration time.Duration) *RLock {
	return &RLock{
		key: key,
		reqID: reqID,
		duration: duration,
	}
}

func (r *RLock) TryLock() bool {
	result, err := RedisDB.SetNX(r.key, r.reqID, r.duration).Result()
	if err != nil {
		log.Errorf("redis lock failed, key=%v, reqID=%v, err=%v", r.key, r.reqID, err.Error())
		return false
	}
	log.Infof("redis lock success=%v, key=%v, reqID=%v",result, r.key, r.reqID)
	return result
}

func (r *RLock) UnLock() error {
	if RedisDB.Get(r.key).Val() != r.reqID {
		return errors.New("other requested locks cannot be unlocked")
	}
	_, err := RedisDB.Del(r.key).Result()
	if err != nil {
		log.Errorf("redis lock unloack failed, key=%v, reqID=%v", r.key, r.reqID)
		return err
	}
	return err
}

func (r *RLock) RenewLock() error {
	_, err := RedisDB.Expire(r.key, r.duration).Result()
	if err != nil {
		log.Errorf("redis lock renew failed, key=%v, reqID=%v", r.key, r.reqID)
		return err
	}
	log.Infof("redis lock renew succeed, key=%v, reqID=%v", r.key, r.reqID)
	return err
}