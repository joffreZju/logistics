package conn

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

//conninfo 的基本结构
type Conn struct {
	Id          string
	Dbtype      string
	Name        string
	Host        string
	Port        int
	DbUser      string
	Passwd      string
	Params      string
	Dbname      string
	Prefix      string
	Db          *gorm.DB
	Status      bool
	Lastusetime time.Time
}
