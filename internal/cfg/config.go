package cfg

import (
	"os"

	"github.com/google/uuid"
	"github.com/iegad/kraken/conf"
	"gopkg.in/yaml.v3"
)

var Instance *config

type config struct {
	Server *conf.Server `json:"server" yaml:"server"`
	Etcd   *conf.Etcd   `json:"etcd"   yaml:"etcd"`
	MySql  *conf.MySql  `json:"mysql"  yaml:"mysql"`
}

func Init(fname string) error {
	data, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	conf := &config{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return err
	}

	if len(conf.Server.ID) == 0 {
		conf.Server.ID = uuid.New().String()
		data, err = yaml.Marshal(conf)
		if err != nil {
			return err
		}

		err = os.WriteFile(fname, data, 0755)
		if err != nil {
			return err
		}
	}

	Instance = conf

	return nil
}
