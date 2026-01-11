package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-konsultin/errk"
	"github.com/konsultin/project-goes-here/internal/svc-core/constant"
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/project-goes-here/pkg/redis"
)

func (r *Repository) FindSessionByXid(xid string) (*model.AuthSession, error) {
	key := fmt.Sprintf("%s%s", constant.RedisSessionPrefix, xid)

	val, err := r.redis.Get(key)
	if redis.IsNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errk.Trace(err)
	}

	var m model.AuthSession
	if err := json.Unmarshal([]byte(val), &m); err != nil {
		return nil, errk.Trace(err)
	}

	return &m, nil
}

func (r *Repository) DeleteSessionByXid(xid string) error {
	key := fmt.Sprintf("%s%s", constant.RedisSessionPrefix, xid)
	_, err := r.redis.Del(key)
	if err != nil {
		return errk.Trace(err)
	}
	return nil
}

func (r *Repository) InsertAuthSession(session *model.AuthSession) error {
	key := fmt.Sprintf("%s%s", constant.RedisSessionPrefix, session.Xid)

	data, err := json.Marshal(session)
	if err != nil {
		return errk.Trace(err)
	}

	lifetime := r.config.UserSessionLifetime
	if lifetime == 0 {
		lifetime = 3600 * time.Second
	}

	err = r.redis.Set(key, data, lifetime)
	if err != nil {
		return errk.Trace(err)
	}

	return nil
}
