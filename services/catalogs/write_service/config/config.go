package config

import (
	"emperror.dev/errors"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gormPostgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	postgres "github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres_pgx"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/probes"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "catalogs write microservice config path")
}

type Config struct {
	DeliveryType     string                         `mapstructure:"deliveryType" env:"DeliveryType"`
	ServiceName      string                         `mapstructure:"serviceName" env:"ServiceName"`
	Logger           *logger.LogConfig              `mapstructure:"logger" envPrefix:"Logger_"`
	GRPC             *grpc.GrpcConfig               `mapstructure:"grpc" envPrefix:"GRPC_"`
	Http             *customEcho.EchoHttpConfig     `mapstructure:"http" envPrefix:"Http_"`
	Context          Context                        `mapstructure:"context" envPrefix:"Context_"`
	Postgresql       *postgres.Config               `mapstructure:"postgres" envPrefix:"Postgresql_"`
	GormPostgres     *gormPostgres.Config           `mapstructure:"gormPostgres" envPrefix:"GormPostgres_"`
	RabbitMQ         *config.RabbitMQConfig         `mapstructure:"rabbitmq" envPrefix:"RabbitMQ_"`
	Probes           probes.Config                  `mapstructure:"probes" envPrefix:"Probes_"`
	Jaeger           *tracing.Config                `mapstructure:"jaeger" envPrefix:"Jaeger_"`
	EventStoreConfig eventstroredb.EventStoreConfig `mapstructure:"eventStoreConfig" envPrefix:"EventStoreConfig_"`
}

type Context struct {
	Timeout int `mapstructure:"timeout" env:"Timeout"`
}

func InitConfig(environment string) (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			//https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
			//https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
			d, err := dirname()
			if err != nil {
				return nil, err
			}

			configPath = d
		}
	}

	cfg := &Config{}

	viper.SetConfigName(fmt.Sprintf("config.%s", environment))
	viper.AddConfigPath(configPath)
	viper.SetConfigType(constants.Json)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.WrapIf(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.WrapIf(err, "viper.Unmarshal")
	}

	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	grpcPort := os.Getenv(constants.GrpcPort)
	if grpcPort != "" {
		cfg.GRPC.Port = grpcPort
	}

	postgresHost := os.Getenv(constants.PostgresqlHost)
	if postgresHost != "" {
		cfg.Postgresql.Host = postgresHost
	}
	postgresPort := os.Getenv(constants.PostgresqlPort)
	if postgresPort != "" {
		cfg.Postgresql.Port = postgresPort
	}
	jaegerAddr := os.Getenv(constants.JaegerHostPort)
	if jaegerAddr != "" {
		cfg.Jaeger.HostPort = jaegerAddr
	}

	return cfg, nil
}

func filename() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filename, nil
}

func dirname() (string, error) {
	filename, err := filename()
	if err != nil {
		return "", err
	}
	return filepath.Dir(filename), nil
}
