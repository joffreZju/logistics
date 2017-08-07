package all

import (
	// Initialize all adapters by importing this package
	_ "github.com/bbjj040471/transporter/adaptor/elasticsearch"
	_ "github.com/bbjj040471/transporter/adaptor/file"
	_ "github.com/bbjj040471/transporter/adaptor/mongodb"
	_ "github.com/bbjj040471/transporter/adaptor/postgres"
	_ "github.com/bbjj040471/transporter/adaptor/rabbitmq"
	_ "github.com/bbjj040471/transporter/adaptor/rethinkdb"
)
