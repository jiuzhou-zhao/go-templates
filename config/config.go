package config

import (
	"sync"

	"github.com/sgostarter/i/l"
	"github.com/sgostarter/libconfig"
	"github.com/sgostarter/liblogrus"
	"github.com/sgostarter/libservicetoolset/servicetoolset"
)

type Config struct {
	Logger l.Wrapper `yaml:"-"`

	GRPCListen4Server    string                            `yaml:"GRPCListen4Server"`
	GRPCTLSConfig4Server *servicetoolset.GRPCTlsFileConfig `yaml:"GRPCTLSConfig4Server"`

	GRPCServerAddress4Client string `yaml:"GRPCServerAddress4Client"`
	GRPCTLSConfig4Client     *servicetoolset.GRPCTlsFileConfig
}

var (
	_cfg  Config
	_once sync.Once
)

func GetConfig() *Config {
	_once.Do(func() {
		_cfg.Logger = l.NewWrapper(liblogrus.NewLogrus())
		_cfg.Logger.GetLogger().SetLevel(l.LevelDebug)

		_, err := libconfig.LoadOnConfigPath("config.yaml", []string{
			"./", "../", "../../", "../../../",
		}, &_cfg)
		if err != nil {
			panic("load config: " + err.Error())
		}
	})

	return &_cfg
}
