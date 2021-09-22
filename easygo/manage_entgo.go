package easygo

import (
	"context"
	"fmt"
	"grpc_gateway/ent"

	"log"

	"entgo.io/ent/dialect"
)

type PostgresCfg struct {
	Host     string
	Port     int
	User     string
	PassWord string
	DbName   string
}

//连接管理
type EntManager struct {
	Psql   *PostgresCfg
	Client *ent.Client
}

func NewEntManager() *EntManager { // services map[string]interface{},
	p := &EntManager{}
	p.Init()
	return p
}

//初始化
func (e *EntManager) Init() {
}

func (e *EntManager) Open(cfg *PostgresCfg) *ent.Client {
	pdqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.PassWord, cfg.DbName)
	client, err := ent.Open(dialect.Postgres, pdqlInfo)
	if err != nil {
		log.Fatalf("failed opening connection to Postgres: %v", err)
	}
	// defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}
