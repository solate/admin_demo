package initialize

import (
	"context"
	"fmt"
	"time"

	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"
)

// InitFactory 初始化工厂数据
func InitFactory(ctx context.Context, tx *generated.Tx, tenantCode string) ([]string, error) {
	// 生成工厂ID
	ids, err := idgen.GenerateUUIDs(3)
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()
	factories := []struct {
		id    string
		name  string
		addr  string
		phone string
	}{
		{ids[0], "北京工厂", "北京市朝阳区", "010-12345678"},
		{ids[1], "上海工厂", "上海市浦东新区", "021-87654321"},
		{ids[2], "深圳工厂", "深圳市龙岗区", "0755-98765432"},
	}

	var factoryIDs []string
	for _, f := range factories {
		_, err := tx.Factory.Create().
			SetCreatedAt(now).
			SetUpdatedAt(now).
			SetTenantCode(tenantCode).
			SetFactoryID(f.id).
			SetFactoryName(f.name).
			SetAddress(f.addr).
			SetContactPhone(f.phone).
			SetStatus(1).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed creating factory: %v", err)
		}
		factoryIDs = append(factoryIDs, f.id)
		fmt.Printf("factory: %s, name: %s, address: %s\n", f.id, f.name, f.addr)
	}

	return factoryIDs, nil
}
