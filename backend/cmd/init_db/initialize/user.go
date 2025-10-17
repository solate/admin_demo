package initialize

import (
	"context"
	"fmt"
	"time"

	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"
	"admin_backend/pkg/utils/passwordgen"
)

const (
	Password = "admin@123"
)

// InitUser 初始化用户数据
func InitUser(ctx context.Context, tx *generated.Tx, tenantCode string) (string, error) {

	// 生成所需的UUID
	ids, err := idgen.GenerateUUIDs(3)
	if err != nil {
		return "", err
	}

	now := time.Now().UnixMilli()
	userID := ids[1]

	salt, err := passwordgen.GenerateSalt()
	if err != nil {
		return "", fmt.Errorf("generate salt failed: %v", err)
	}

	hashedPassword, err := passwordgen.Argon2Hash(Password, salt)
	if err != nil {
		return "", fmt.Errorf("hash password failed: %v", err)
	}

	user, err := tx.SysUser.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(tenantCode).
		SetUserID(userID).
		SetPhone("18888888888").
		SetUserName("admin").
		SetPwdHashed(hashedPassword).
		SetPwdSalt(salt).
		SetStatus(1).
		SetName("admin").
		SetEmail("admin123@qq.com").
		SetSex(1).
		Save(ctx)

	if err != nil {
		return "", fmt.Errorf("failed creating schema user resources: %v", err)
	}

	fmt.Printf("user: %s, password: %s\n", user.UserName, Password)

	return userID, nil
}
