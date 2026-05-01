package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute * 10

func (s *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	key := fmt.Sprintf("user:%d", userID)
	fmt.Println("ek tha")

	val, err := s.rdb.Get(ctx, key).Result()

	fmt.Println("cccc")
	if err != nil {
		fmt.Println(err, "1")
		if err == redis.Nil {
			fmt.Println(err, "2")
			return nil, nil
		}
		return nil, err
	}
	var user store.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		fmt.Println(err, "3")
		return nil, err
	}
	return &user, nil

}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	fmt.Println("chrrrrrrrrrrrrrrrrrr")
	key := fmt.Sprintf("user:%d", user.ID)
	_json, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = s.rdb.Set(ctx, key, _json, UserExpTime).Err()

	return err
}

func (s *UserStore) Delete(ctx context.Context, userID int64) error {

	key := fmt.Sprintf("user:%d", userID)
	return s.rdb.Del(ctx, key).Err()

}
