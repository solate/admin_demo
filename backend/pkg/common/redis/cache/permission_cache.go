package cache

import (
	"fmt"
	"time"
)

const (
	// PermissionChildrenKey 权限子节点缓存key
	PermissionChildrenKey = "user:perm:%s"
	// PermissionCacheExpiration 权限缓存过期时间
	PermissionCacheExpiration = 24 * time.Hour
)

// CachePermissionChildren 缓存所有子权限ID
func (r *CacheManager) CachePermissionChildren(permissionId string, ids []string) error {
	key := fmt.Sprintf(PermissionChildrenKey, permissionId)
	return r.Set(key, ids, PermissionCacheExpiration)
}

// GetPermissionChildrenFromCache 从缓存获取子权限ID
func (r *CacheManager) GetPermissionChildrenFromCache(permissionId string) ([]string, error) {
	key := fmt.Sprintf(PermissionChildrenKey, permissionId)
	var ids []string

	// 检查缓存是否存在
	exists, err := r.Exists(key)
	if err != nil {
		return nil, fmt.Errorf("check cache exists failed: %v", err)
	}

	if !exists {
		return nil, fmt.Errorf("cache not found for permission %s", permissionId)
	}

	// 获取缓存数据
	if err := r.Get(key, &ids); err != nil {
		return nil, fmt.Errorf("get cache failed: %v", err)
	}

	return ids, nil
}
