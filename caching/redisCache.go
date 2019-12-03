package caching

import (
	"net"
	"time"
	"github.com/go-redis/redis"
)

var (
	address string = "localhost:6379"
	password string = "" // No password
	db int = 0
)

// RedisCache is a cacher that used redis DB to manage the data.
type RedisCache struct{
	redisClient *redis.Client
}

// Creates new RedisCache object.
func NewRedisCache() *RedisCache{
	return &RedisCache{}
}

// Initialize new redis cache object.
func (rc *RedisCache) InitializeCache() error{
	rdb := redis.NewClient(&redis.Options{
	Addr:     address,
	Password: password,
	DB:       db,
	})

	_, err := rdb.Ping().Result()
	rc.redisClient = rdb
	return err
}

// Find the related ip addresses.
func (rc *RedisCache) GetIPS(domain string) ([]net.IP, error){
	var ips []net.IP
	results, err := rc.redisClient.SMembers(domain).Result()
	if err != nil {
		return ips, err
	}

	for _, result := range results{
		ips = append(ips, net.ParseIP(result))
	}
	return ips, err

}

// Updates the assocciated ip address for the given domain.
// If the domain already exists in the db, this function will override.
func (rc *RedisCache) UpdateDomain(domain string, ips []net.IP, expiration time.Duration)(error){
	var err error
	var ipString []string
	_, err = rc.redisClient.Del(domain).Result()
	if err != nil{
		return err
	}

	for _, ip := range ips{
		ipString = append(ipString, ip.String())
	}
	_, err = rc.redisClient.SAdd(domain, ipString).Result()
	if err != nil{
		return err
	}
	_, err = rc.redisClient.Expire(domain, expiration).Result()
	return err

}
