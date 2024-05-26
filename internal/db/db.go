package db

import (
	"database/sql"
	"fmt"
	"hash/fnv"

	"github.com/cjnghn/db-shard-example/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

type DBShard struct {
	DB   *sql.DB
	Name string
}

var Shards []DBShard

func GetShardByUUID(uuid string) DBShard {
	h := fnv.New32a()
	h.Write([]byte(uuid))
	shardIndex := h.Sum32() % uint32(len(Shards))

	return Shards[shardIndex]
}

const userTableDDL = `
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
`

func InitShards(config *config.Config) error {
	var shards []DBShard
	for _, shardConfig := range config.DBShards {
		shard, err := initShard(shardConfig)
		if err != nil {
			return err
		}
		shards = append(shards, *shard)
	}
	Shards = shards
	return nil
}

func initShard(shardConfig config.DBConfig) (*DBShard, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		shardConfig.User,
		shardConfig.Password,
		shardConfig.Host,
		shardConfig.Port,
		shardConfig.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", shardConfig.Name, err)
	}

	if err := createUserTable(db, shardConfig.Name); err != nil {
		db.Close()
		return nil, err
	}
	return &DBShard{DB: db, Name: shardConfig.Name}, nil
}

func createUserTable(db *sql.DB, shardName string) error {
	_, err := db.Exec(userTableDDL)
	if err != nil {
		return fmt.Errorf("failed to create user table in %s: %v", shardName, err)
	}
	return nil
}
