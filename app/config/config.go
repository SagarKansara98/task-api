package config

import (
	"fmt"

	"github.com/ardanlabs/conf/v3"
	"github.com/pkg/errors"
)

type Config struct {
	Dns             string `conf:"default:postgresql://host1:123,host2:456/somedb?target_session_attrs=any&application_name=myapp,env:PG_CONNECTION"`
	Port            string `conf:"default:3000"`
	LogLevel        string `conf:"default:info"`
	JwtTokenSeceret string `conf:"default:123456789"`
}

func ParseConfig() (Config, error) {
	var cfg Config
	_, err := conf.Parse("", &cfg)
	if err != nil {
		return cfg, errors.Wrap(err, "ParseConfig: parsing config")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return cfg, errors.Wrap(err, "ParseConfig: parsing config")
	}

	fmt.Printf("------------------- App Config -------------------\n %+v \n --------------------------------------------- \n", out)

	return cfg, nil
}
