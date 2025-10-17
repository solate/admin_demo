package casbin

import (
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"context"
	"testing"
)

func setupTestDB(_ *testing.T) (*ent.Client, error) {
	dataSource := "user=root password=root host=127.0.0.1 port=5432 dbname=testdb sslmode=disable"
	client, err := ent.NewClient(context.Background(), dataSource)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func cleanupTest(client *generated.Client) {
	if client != nil {
		client.Close()
	}
}

func TestNewCasbin(t *testing.T) {
	client, err := setupTestDB(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanupTest(client.Client)

	e, err := NewCasbin(client)
	if err != nil {
		t.Error(err)
		return
	}

	// 测试添加策略
	sub := "18811111111"
	tenant := "demo"
	obj := "Dashboard"
	act := "read"
	tType := "menu"

	// 添加策略
	ok, err := e.AddPolicy(sub, tenant, obj, act, tType)
	if err != nil {
		t.Error(err)
		return
	}
	if !ok {
		t.Error("Failed to add policy")
		return
	}

	// 测试权限验证
	ok, err = e.Enforce(sub, tenant, obj, act, tType)
	if err != nil {
		t.Error(err)
		return
	}
	if !ok {
		t.Error("Expected permission to be granted")
	}

	// 测试删除策略
	ok, err = e.RemovePolicy(sub, tenant, obj, act, tType)
	if err != nil {
		t.Error(err)
		return
	}
	if !ok {
		t.Error("Failed to remove policy")
		return
	}

	// 验证策略已被删除
	ok, err = e.Enforce(sub, tenant, obj, act, tType)
	if err != nil {
		t.Error(err)
		return
	}
	if ok {
		t.Error("Expected permission to be denied after removal")
	}
}
