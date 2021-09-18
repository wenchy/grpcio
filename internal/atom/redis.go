package atom

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

func (atom *Atom) InitRedis(addrs []string, password string) error {
	// client := redis.NewClient(&redis.Options{
	// 	Addr:     address,
	// 	Password: password, // password set
	// 	DB:       0,        // use default DB
	// })

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
		// To route commands by latency or randomly, enable one of the following.
		//RouteByLatency: true,
		//RouteRandomly: true,
	})

	var ctx = context.Background()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("redis ping-pong: %s, err: %v", pong, err)
		return err
	}

	atom.RedisClient = client
	return nil
}
