package saga

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
)

type RedisCarrier struct {
	Active  bool
	Options *RedisCarrierConfig
	Client  *redis.Client
}

func getRedisCarrierInstance() *RedisCarrier {
	return &RedisCarrier{
		Active:  false,
		Options: &RedisCarrierConfig{},
	}
}

func (rc *RedisCarrier) IsActive() bool {
	return rc.Active
}

func (rc *RedisCarrier) Activate() error {
	if rc.Client == nil {
		return errors.New("Client not set")
	}
	rc.Active = true
	return nil
}

func (rc *RedisCarrier) SetOptions(opts CarrierConfig) error {
	cfg, _ := opts.(*RedisCarrierConfig)
	err := cfg.Verify()
	if err != nil {
		return err
	}
	rc.Options = cfg
	err = rc.connect()
	if err != nil {
		return err
	}
	return nil
}

// connect establishes a connection to the Redis server.
func (rc *RedisCarrier) connect() error {
	rc.Client = redis.NewClient(&redis.Options{
		// Address for connecting to the Redis server
		Addr: rc.Options.Addr,
		// Leave blank if no password is set
		Password: rc.Options.Password,
		// Set to 0 to use the default database.
		DB: 0,
	})

	// rc.Client = asynq.NewClient(asynq.RedisClientOpt{
	// 	Addr:     rc.Options.Addr,
	// 	Password: rc.Options.Password,
	// })
	return nil
}

func (rc *RedisCarrier) disconnect() error {
	if rc.Client != nil {
		rc.Client.Close()
	}
	return nil
}

func (rc *RedisCarrier) Push(Message string, Data interface{}) error {
	dataJSON, err := json.Marshal(Data)
	if err != nil {
		return err
	}
	rc.Client.Publish(context.Background(), Message, dataJSON)
	return nil
}
func (rc *RedisCarrier) AddEventsListener(func(Message string, Data interface{})) error {
	return nil
}

// ************************************************************
/* Redis Carrier Configuration */
type RedisCarrierConfig struct {
	Addr     string
	Password string
}

func (redisCfg *RedisCarrierConfig) Verify() error {
	return nil
}
