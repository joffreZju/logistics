package all

import (
	// ensures init functions get called
	_ "github.com/bbjj040471/transporter/adaptor/elasticsearch/clients/v1"
	_ "github.com/bbjj040471/transporter/adaptor/elasticsearch/clients/v2"
	_ "github.com/bbjj040471/transporter/adaptor/elasticsearch/clients/v5"
)
