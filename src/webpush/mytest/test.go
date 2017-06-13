package main

import (
	"allsum_oa/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

func Sendxxx(m *model.Message) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(time.Second * 10)
				c, err := net.DialTimeout(netw, addr, 5*time.Second)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	webpush := "http://127.0.0.1:8009/newmsg"
	body, _ := json.Marshal(m)
	fmt.Printf("====%s\n", string(body))

	request, err := http.NewRequest("POST", webpush, bytes.NewReader(body))
	if err != nil {
		fmt.Println("http.NewRequest:", err)
		return
	}
	_, err = client.Do(request)
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	xxx := new(model.Message)
	xxx.Title = "审批通知"
	xxx.Content = "请假单:xxxx需要您审批"
	xxx.UserId = 10
	Sendxxx(xxx)
}
