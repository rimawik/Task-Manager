package cache

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/task-manager/config"
	"gopkg.in/redis.v5"
)

type Rdb struct {
	RDBClient *redis.Client
}

func ConnectToRedis(cfg config.Config) (*Rdb, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Errorf("couldn't connect to redis: %v", err)
		return &Rdb{}, err
	}
	return &Rdb{RDBClient: rdb}, nil

}

func (rdb *Rdb) SetJson(key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		log.Errorf("couldn't unmarshal json value :%v", err)
		return err
	}
	err2 := rdb.RDBClient.Set(key, string(bytes), 3600*time.Second)
	if err2 != nil {
		log.Errorf("Could not set key: %v", err2)
		return err2.Err()
	}
	return nil
}

func (rdb *Rdb) Set(key string, value string) error {

	err := rdb.RDBClient.Set(key, value, time.Hour)
	if err != nil {
		log.Errorf("Could not insert key: %v", err)
		return err.Err()
	}
	return nil
}
func (rdb *Rdb) Get(key string) (value string, err error) {

	value, err = rdb.RDBClient.Get(key).Result()
	if err != nil {
		log.Errorf("Could not get key: %v", err)
		return "", err
	}
	return value, nil

}

func (rdb *Rdb) Del(key string) (err error) {

	re := rdb.RDBClient.Del(key)
	if re.Err() != nil {
		log.Errorf("Could not get key: %v", err)
		return err
	}
	return nil

}
