package cache

import (
	"common/lib/keycrypt"
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
)

var (
	Cache      cache.Cache
	LocalCache cache.Cache
)

func Init(key string) (err error) {
	LocalCache = cache.NewMemoryCache()
	conf := beego.AppConfig.String("cache::params")
	if len(key) > 0 {
		m := map[string]string{}
		err = json.Unmarshal([]byte(conf), &m)
		if err != nil {
			return
		}
		if pass, ok := m["password"]; ok && len(pass) > 0 {
			pass, err = keycrypt.Decode(key, pass)
			if err != nil {
				return
			}
			m["password"] = pass
			c, err := json.Marshal(m)
			if err != nil {
				return err
			}
			conf = string(c)
		}
	}
	Cache, err = cache.NewCache("redis", conf)
	return
}
