package config

// DBConfig holds the configuration for a single database shard.
type DBConfig struct {
    Name     string
    User     string
    Password string
    Host     string
    Port     string
}

type Config struct {
    DBShards []DBConfig
}

func GetConfig() *Config {
    return &Config{
        DBShards: []DBConfig{
            {
                Name:     "shard1",
                User:     "root",
                Password: "password",
                Host:     "localhost",
                Port:     "3307",
            },
            {
                Name:     "shard2",
                User:     "root",
                Password: "password",
                Host:     "localhost",
                Port:     "3308",
            },
            // Add more shards as needed
        },
    }
}