package cfg

import (
	"os"

	"github.com/iegad/kraken/conf"
	"github.com/iegad/kraken/utils"
	"gopkg.in/yaml.v3"
)

var (
	Instance *config
)

type config struct {
	Server *conf.Server `json:"server" yaml:"server"`
	Etcd   *conf.Etcd   `json:"etcd"   yaml:"etcd"`
	MySql  *conf.MySql  `json:"mysql"  yaml:"mysql"`
	Redis  *conf.Redis  `json:"redis"  yaml:"redis"`
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
		conf.Server.ID = utils.UUID_String()
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
