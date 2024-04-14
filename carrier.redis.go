package saga

type RedisCarrier struct {
	IsInit  bool
	Active  bool
	Options CarrierConfig
}

type RedisCarrierOption struct {
}

func getRedisCarrierInstance() *RedisCarrier {
	return &RedisCarrier{
		IsInit:  false,
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

func (redisOpt *RedisCarrierOption) Verify() error {
	return nil
}
