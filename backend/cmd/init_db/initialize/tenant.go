package initialize

import (
	"context"
	"fmt"
	"time"

	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"
)

const (
	TenantCode = "default"
)

// InitTenant 初始化租户数据
func InitTenant(ctx context.Context, tx *generated.Tx) (string, error) {

	// 生成所需的UUID
	ids, err := idgen.GenerateUUIDs(3)
	if err != nil {
		return "", err
	}

	now := time.Now().UnixMilli()
	tenantID := ids[0]
	_, err = tx.Tenant.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantID(tenantID).
		SetCode(TenantCode).
		SetName("默认租户").
		SetDescription("默认租户").
		SetStatus(1).
		Save(ctx)
	if err != nil {
		return "", fmt.Errorf("failed creating schema tenant resources: %v", err)
	}

	fmt.Printf("tenant: %s, code: %s, name: %s, description: %s, status: %d\n",
		tenantID, TenantCode, "默认租户", "默认租户", 1)

	return tenantID, nil
}
