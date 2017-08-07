package file

import (
	"os"

	"github.com/bbjj040471/transporter/client"
)

// Session serves as a wrapper for the underlying file
type Session struct {
	file *os.File
}

var _ client.Session = &Session{}
