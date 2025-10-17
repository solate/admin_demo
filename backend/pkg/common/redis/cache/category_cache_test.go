package cache

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestCacheCategoryChildren(t *testing.T) {
// 	// 初始化测试用的Redis客户端
// 	testRedis := NewRedis(RedisConfig{
// 		Host: "127.0.0.1",
// 		Port: 6379,
// 		Pass: "123456",
// 		Type: "node",
// 	})
// 	cacheManager := NewCacheManager(testRedis)

// 	tests := []struct {
// 		name       string
// 		categoryId string
// 		ids        []string
// 		wantErr    bool
// 	}{
// 		{
// 			name:       "cache category children successfully",
// 			categoryId: "1",
// 			ids:        []string{"2", "3", "4"},
// 			wantErr:    false,
// 		},
// 		{
// 			name:       "cache empty children list",
// 			categoryId: "5",
// 			ids:        []string{},
// 			wantErr:    false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := cacheManager.CacheCategoryChildren(tt.categoryId, tt.ids)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func TestGetCategoryChildrenFromCache(t *testing.T) {
// 	// 初始化测试用的Redis客户端
// 	testRedis := NewRedis(RedisConfig{
// 		Host: "127.0.0.1",
// 		Port: 6379,
// 		Pass: "123456",
// 		Type: "node",
// 	})
// 	cacheManager := NewCacheManager(testRedis)

// 	// 预先设置测试数据
// 	testCategoryId := "1"
// 	testIds := []string{"2", "3", "4"}
// 	err := cacheManager.CacheCategoryChildren(testCategoryId, testIds)
// 	assert.NoError(t, err)

// 	tests := []struct {
// 		name       string
// 		categoryId string
// 		wantIds    []string
// 		wantErr    bool
// 	}{
// 		{
// 			name:       "get existing category children",
// 			categoryId: testCategoryId,
// 			wantIds:    testIds,
// 			wantErr:    false,
// 		},
// 		{
// 			name:       "get non-existing category children",
// 			categoryId: "999",
// 			wantIds:    nil,
// 			wantErr:    true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ids, err := cacheManager.GetCategoryChildrenFromCache(tt.categoryId)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.wantIds, ids)
// 			}
// 		})
// 	}
// }

// func TestGetCategoryAndChildrenIds(t *testing.T) {
// 	// 初始化测试用的Redis客户端
// 	testRedis := NewRedis(RedisConfig{
// 		Host: "127.0.0.1",
// 		Port: 6379,
// 		Pass: "123456",
// 		Type: "node",
// 	})
// 	cacheManager := NewCacheManager(testRedis)

// 	// 预先设置测试数据
// 	testCategoryId := "1"
// 	testChildrenIds := []string{"2", "3", "4"}
// 	err := cacheManager.CacheCategoryChildren(testCategoryId, testChildrenIds)
// 	assert.NoError(t, err)

// 	tests := []struct {
// 		name       string
// 		categoryId string
// 		wantIds    []string
// 		wantErr    bool
// 	}{
// 		{
// 			name:       "get category and children ids successfully",
// 			categoryId: testCategoryId,
// 			wantIds:    append(testChildrenIds, testCategoryId),
// 			wantErr:    false,
// 		},
// 		{
// 			name:       "get non-existing category and children ids",
// 			categoryId: "999",
// 			wantIds:    nil,
// 			wantErr:    true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ids, err := cacheManager.GetCategoryAndChildrenIds(tt.categoryId)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.wantIds, ids)
// 			}
// 		})
// 	}
// }
