package saga

type RedisCarrier struct {
	Active  bool
	Options CarrierConfig
}

func getRedisCarrierInstance() *RedisCarrier {
	return &RedisCarrier{
		Active:  false,
		Options: &RedisCarrierOption{},
	}
}

func (rc *RedisCarrier) IsActive() bool {
	return rc.Active
}
func (rc *RedisCarrier) SetOptions(opts CarrierConfig) error {
	val, _ := opts.(*InMemoryCarrierConfig)
	err := val.Verify()
	if err != nil {
		return err
	}
	rc.Options = opts
	return nil
}

func (mem *RedisCarrier) Push(Message string, Data interface{}) error {
	return nil
}
func (mem *RedisCarrier) AddListener(func(Message string, Data interface{})) error {
	return nil
}

// ************************************************************
/* Redis Carrier Configuration */
type RedisCarrierOption struct {
}

func (redisOpt *RedisCarrierOption) Verify() error {
	return nil
}
