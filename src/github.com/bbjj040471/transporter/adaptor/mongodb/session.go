package mongodb

import (
	"github.com/bbjj040471/transporter/client"
	mgo "gopkg.in/mgo.v2"
)

// Session serves as a wrapper for the underlying mgo.Session
type Session struct {
	mgoSession *mgo.Session
}

var _ client.Session = &Session{}

// Close implements necessary calls to cleanup the underlying mgo.Session
func (s *Session) Close() {
	s.mgoSession.Close()
}
