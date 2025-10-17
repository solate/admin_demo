package cache

import (
	"time"
)

const (
	// CategoryKey 分类树缓存key
	CategoryKey = "category:%s"

	// CategoryChildrenKey 分类子节点缓存key
	CategoryChildrenKey = "category:children:%s"
	// CategoryCacheExpiration 分类缓存过期时间
	CategoryCacheExpiration = 24 * time.Hour
)

// // CacheCategoryChildren 缓存所有子分类ID
// func (r *CacheManager) CacheCategoryChildren(categoryId string, ids []string) error {
// 	key := fmt.Sprintf(CategoryChildrenKey, categoryId)
// 	return r.Set(key, ids, CategoryCacheExpiration)
// }

// // GetCategoryChildrenFromCache 从缓存获取子分类ID
// func (r *CacheManager) GetCategoryChildrenFromCache(categoryId string) ([]string, error) {
// 	key := fmt.Sprintf(CategoryChildrenKey, categoryId)
// 	var ids []string

// 	// 检查缓存是否存在
// 	exists, err := r.Exists(key)
// 	if err != nil {
// 		return nil, fmt.Errorf("check cache exists failed: %v", err)
// 	}

// 	if !exists {
// 		return nil, fmt.Errorf("cache not found for category %s", categoryId)
// 	}

// 	// 获取缓存数据
// 	if err := r.Get(key, &ids); err != nil {
// 		return nil, fmt.Errorf("get cache failed: %v", err)
// 	}

// 	return ids, nil
// }

// // 获取分类及其所有子分类的ID
// func (r *CacheManager) GetCategoryAndChildrenIds(categoryId string) ([]string, error) {
// 	// 从缓存中获取子分类ID
// 	childrenIds, err := r.GetCategoryChildrenFromCache(categoryId)
// 	if err != nil {
// 		return nil, fmt.Errorf("get category children from cache failed: %v", err)
// 	}

// 	// 将子分类ID添加到结果列表中
// 	result := append(childrenIds, categoryId)
// 	return result, nil
// }
