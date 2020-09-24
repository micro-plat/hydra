package redis

func (r *redisRegistry) Exists(path string) (bool, error) {
	rpath := joinR(path)
	rs, err := r.client.Exists(rpath).Result()
	if err == nil && rs == 1 {
		return true, nil
	}

	return false, err
}
