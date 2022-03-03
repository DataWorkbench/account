package config

import (
	"fmt"
	"github.com/DataWorkbench/common/lib/iaas"
	"io/ioutil"
	"os"
	"time"

	"github.com/DataWorkbench/common/gormwrap"
	"github.com/DataWorkbench/common/grpcwrap"
	"github.com/DataWorkbench/common/gtrace"
	"github.com/DataWorkbench/common/metrics"
	"github.com/DataWorkbench/common/rediswrap"
	"github.com/DataWorkbench/common/utils/logutil"
	"github.com/DataWorkbench/loader"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

var FilePath string

const (
	envPrefix = "ACCOUNT"
)

type Config struct {
	LogConfig     *logutil.Config        `json:"log" yaml:"log" env:"LOG,default=" validate:"required"`
	GRPCLog       *grpcwrap.LogConfig    `json:"grpc_log"       yaml:"grpc_log"       env:"GRPC_LOG"            validate:"required"`
	GRPCServer    *grpcwrap.ServerConfig `json:"grpc_server"    yaml:"grpc_server"    env:"GRPC_SERVER"         validate:"required"`
	MetricsServer *metrics.Config        `json:"metrics_server" yaml:"metrics_server" env:"METRICS_SERVER"      validate:"required"`
	MySQL         *gormwrap.MySQLConfig  `json:"mysql"          yaml:"mysql"          env:"MYSQL"               validate:"required"`
	Redis         *rediswrap.RedisConfig `json:"redis"          yaml:"redis"          env:"REDIS"               validate:"required"`
	Qingcloud     *iaas.Config           `json:"iaas"           yaml:"iaas"           env:"IAAS"                validate:"required"`
	Source        string                 `json:"source"         yaml:"source"         env:"SOURCE"              validate:"required"`
	Tracer        *gtrace.Config         `json:"tracer"         yaml:"tracer"         env:"TRACER"              validate:"required"`
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
