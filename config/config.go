package config

import (
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type (
	AppConfig struct {
		Server         Server
		DatabaseServer DatabaseServer
	}

	DatabaseServer struct {
		Hostname string
		Username string
		Password string
	}

	Server struct {
		Address string
		Port    string
	}
)

func New(file string) (*AppConfig, error) {
	if file != "" {
		viper.SetConfigType(file)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("ReadInConfig", slog.Any("err", err))
		slog.Info("Load configuration in system environment.")
	}

	return &AppConfig{
		Server:         getServer(),
		DatabaseServer: getDatabaseServer(),
	}, nil
}

func getServer() Server {
	return Server{
		Address: viper.GetString("address"),
		Port:    viper.GetString("port"),
	}
}

func getDatabaseServer() DatabaseServer {
	return DatabaseServer{
		Username: viper.GetString("mongo.username"),
		Password: viper.GetString("mongo.password"),
		Hostname: viper.GetString("mongo.hostname"),
	}
}
