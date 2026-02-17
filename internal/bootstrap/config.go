package bootstrap

import (
	"github.com/spf13/viper"
)

type Config struct {
	AppPort  string
	GRPCPort string
	DBUrl    string
	Redis    RedisConfig
	Kafka    KafkaConfig
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type KafkaConfig struct {
	Brokers  []string
	GroupID  string
	Topics   KafkaTopics
}

type KafkaTopics struct {
	UserEvents string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		AppPort:  viper.GetString("app.port"),
		GRPCPort: viper.GetString("app.grpc_port"),
		DBUrl:    viper.GetString("database.url"),
		Redis: RedisConfig{
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
		Kafka: KafkaConfig{
			Brokers: viper.GetStringSlice("kafka.brokers"),
			GroupID: viper.GetString("kafka.group_id"),
			Topics: KafkaTopics{
				UserEvents: viper.GetString("kafka.topics.user_events"),
			},
		},
	}, nil
}
