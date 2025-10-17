package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

var testRedis *redis.Redis

func TestMain(m *testing.M) {
	testRedis = NewRedis(RedisConfig{
		Host: "127.0.0.1",
		Port: 6379,
		Pass: "123456",
		Type: "node",
	})
	m.Run()
}

func TestNewRedis(t *testing.T) {
	tests := []struct {
		name string
		cfg  RedisConfig
	}{
		{
			name: "init with default settings",
			cfg: RedisConfig{
				Host: "127.0.0.1",
				Port: 6379,
			},
		},
		{
			name: "init with password and node type",
			cfg: RedisConfig{
				Host: "127.0.0.1",
				Port: 6379,
				Pass: "123456",
				Type: "node",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewRedis(tt.cfg)
			assert.NotNil(t, client)
		})
	}
}

func TestRedis_Set(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		value      string
		expiration int
		wantErr    bool
	}{
		{
			name:       "set string value",
			key:        "test_key",
			value:      "test_value",
			expiration: 100,
			wantErr:    false,
		},
		{
			name:       "set without expiration",
			key:        "test_key_2",
			value:      "test_value_2",
			expiration: 0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.expiration > 0 {
				err = testRedis.Setex(tt.key, tt.value, tt.expiration)
			} else {
				err = testRedis.Set(tt.key, tt.value)
			}

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRedis_Get(t *testing.T) {
	// 预先设置测试数据
	testKey := "test_key"
	testValue := "test_value"
	err := testRedis.Set(testKey, testValue)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		key     string
		wantVal string
		wantErr bool
	}{
		{
			name:    "get existing key",
			key:     testKey,
			wantVal: testValue,
			wantErr: false,
		},
		{
			name:    "get non-existing key",
			key:     "non_existing_key",
			wantVal: "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := testRedis.Get(tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantVal, val)
			}
		})
	}
}

func TestRedis_Del(t *testing.T) {
	// 预先设置测试数据
	testKey := "test_key"
	err := testRedis.Set(testKey, "test_value")
	assert.NoError(t, err)

	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "delete existing key",
			key:     testKey,
			wantErr: false,
		},
		{
			name:    "delete non-existing key",
			key:     "non_existing_key",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testRedis.Del(tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				exists, _ := testRedis.Exists(tt.key)
				assert.False(t, exists)
			}
		})
	}
}

func TestRedis_Exists(t *testing.T) {
	// 预先设置测试数据
	testKey := "test_key"
	err := testRedis.Set(testKey, "test_value")
	assert.NoError(t, err)

	tests := []struct {
		name      string
		key       string
		wantExist bool
		wantErr   bool
	}{
		{
			name:      "key exists",
			key:       testKey,
			wantExist: true,
			wantErr:   false,
		},
		{
			name:      "key does not exist",
			key:       "non_existing_key",
			wantExist: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := testRedis.Exists(tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantExist, exists)
			}
		})
	}
}

func TestRedis_Ttl(t *testing.T) {
	// 预先设置带过期时间的测试数据
	testKey := "test_key"
	err := testRedis.Setex(testKey, "test_value", 100)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		key     string
		wantTTL time.Duration
		wantErr bool
	}{
		{
			name:    "get TTL for existing key",
			key:     testKey,
			wantTTL: 100 * time.Second,
			wantErr: false,
		},
		{
			name:    "get TTL for non-existing key",
			key:     "non_existing_key",
			wantTTL: -2 * time.Second,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ttl, err := testRedis.Ttl(tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.wantTTL > 0 {
					assert.True(t, time.Duration(ttl)*time.Second > 0 && time.Duration(ttl)*time.Second <= tt.wantTTL)
				} else {
					assert.Equal(t, int(tt.wantTTL.Seconds()), ttl)
				}
			}
		})
	}
}
