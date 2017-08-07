package pick

import (
	"github.com/bbjj040471/transporter/function"
	"github.com/bbjj040471/transporter/log"
	"github.com/bbjj040471/transporter/message"
)

func init() {
	function.Add(
		"pick",
		func() function.Function {
			return &Picker{}
		},
	)
}

type Picker struct {
	Fields []string `json:"fields"`
}

func (p *Picker) Apply(msg message.Msg) (message.Msg, error) {
	log.With("msg", msg).Debugln("picking...")
	pluckedMsg := map[string]interface{}{}
	for _, k := range p.Fields {
		if v, ok := msg.Data().AsMap()[k]; ok {
			pluckedMsg[k] = v
		}
	}
	log.With("msg", pluckedMsg).Debugln("...picked")
	return message.From(msg.OP(), msg.Namespace(), pluckedMsg), nil
}
