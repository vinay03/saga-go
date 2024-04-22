package saga

import (
	"github.com/hibiken/asynq"
)

type RedisCarrier struct {
	Active  bool
	Options *RedisCarrierConfig
	Client  *asynq.Client
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

func (rc *RedisCarrier) connect() error {
	rc.Client = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     rc.Options.Host,
		Password: rc.Options.Password,
	})
	return nil
}
func (rc *RedisCarrier) disconnect() error {
	if rc.Client != nil {
		rc.Client.Close()
	}
	return nil
}

func (mem *RedisCarrier) Push(Message string, Data interface{}) error {
	return nil
}
func (mem *RedisCarrier) AddEventsListener(func(Message string, Data interface{})) error {
	return nil
}

// ************************************************************
/* Redis Carrier Configuration */
type RedisCarrierConfig struct {
	Host     string
	Password string
}

func (redisCfg *RedisCarrierConfig) Verify() error {
	return nil
}
