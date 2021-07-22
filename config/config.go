package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/DataWorkbench/common/gormwrap"
	"github.com/DataWorkbench/common/grpcwrap"
	"github.com/DataWorkbench/common/metrics"
	"github.com/DataWorkbench/common/rediswrap"
	"github.com/DataWorkbench/loader"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

var FilePath string

const (
	envPrefix = "ACCOUNT"
)

type QingcloudConfig struct {
	AccessKeyID     string `json:"access_key_id"     yaml:"access_key_id"     env:"ACCESS_KEY_ID"     validate:"required"`
	SecretAccessKey string `json:"secret_access_key" yaml:"secret_access_key" env:"SECRET_ACCESS_KEY" validate:"required"`
	Host            string `json:"host" yaml:"host" env:"HOST" validate:"required"`
	Port            int    `json:"port" yaml:"port" env:"PORT" validate:"required"`
	Protocol        string `json:"protocol" yaml:"protocol" env:"PROTOCOL" validate:"required"`
	Uri             string `json:"uri" yaml:"uri" env:"URI" validate:"required"`
}

type Config struct {
	LogLevel        int8                   `json:"log_level"      yaml:"log_level"      env:"LOG_LEVEL,default=1" validate:"gte=1,lte=5"`
	GRPCLog         *grpcwrap.LogConfig    `json:"grpc_log"       yaml:"grpc_log"       env:"GRPC_LOG"            validate:"required"`
	GRPCServer      *grpcwrap.ServerConfig `json:"grpc_server"    yaml:"grpc_server"    env:"GRPC_SERVER"         validate:"required"`
	MetricsServer   *metrics.Config        `json:"metrics_server" yaml:"metrics_server" env:"METRICS_SERVER"      validate:"required"`
	MySQL           *gormwrap.MySQLConfig  `json:"mysql"          yaml:"mysql"          env:"MYSQL"               validate:"required"`
	Redis           *rediswrap.RedisConfig `json:"redis"          yaml:"redis"          env:"REDIS"               validate:"required"`
	QingcloudConfig *QingcloudConfig       `json:"qingcloud_config"      yaml:"qingcloud_config"      env:"QINGCLOUD_CONFIG"           validate:"required"`
}

func loadFromFile(cfg *Config) (err error) {
	if FilePath == "" {
		return
	}

	fmt.Printf("%s load config from file <%s>\n", time.Now().Format(time.RFC3339Nano), FilePath)

	var b []byte
	b, err = ioutil.ReadFile(FilePath)
	if err != nil && os.IsNotExist(err) {
		return
	}

	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		fmt.Println("parse config file error:", err)
	}

	return
}

// LoadConfig load all configuration from specified file
// Must be set `FilePath` before called
func Load() (cfg *Config, err error) {
	cfg = &Config{}

	_ = loadFromFile(cfg)

	l := loader.New(
		loader.WithPrefix(envPrefix),
		loader.WithTagName("env"),
		loader.WithOverride(true),
	)
	if err = l.Load(cfg); err != nil {
		return
	}

	// output the config content
	fmt.Printf("%s pid=%d the latest configuration: \n", time.Now().Format(time.RFC3339Nano), os.Getpid())
	b, _ := yaml.Marshal(cfg)
	fmt.Println(string(b))

	validate := validator.New()
	if err = validate.Struct(cfg); err != nil {
		return
	}
	return
}
