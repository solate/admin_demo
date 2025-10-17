package ent

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	dataSource := "user=root password=root host=127.0.0.1 port=5432 dbname=testdb sslmode=disable"
	client, err := NewClient(context.Background(), dataSource)
	if err != nil {
		panic(err)
	}

	t.Log("client", client)
}
