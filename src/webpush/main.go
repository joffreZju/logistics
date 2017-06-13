package main

import (
	"allsum_oa/model"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

var Push *PushPool

type Client struct {
	Conn   *websocket.Conn
	UserId int
	InfoC  chan []byte
}

func NewClient(c *websocket.Conn, id int) *Client {
	cl := new(Client)
	cl.Conn = c
	cl.InfoC = make(chan []byte)
	cl.UserId = id
	Push.AddUser(id, &cl.InfoC)
	return cl
}

func (c *Client) Close() {
	c.Conn.Close()
	close(c.InfoC)
}

func NewMsg(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("http read:", err)
		return
	}
	m := new(model.Message)
	err = json.Unmarshal(body, m)
	if err != nil {
		log.Print("json unmarshal:", err)
		return
	}
	if c, ok := Push.GetUserChan(m.UserId); ok {
		*c <- []byte(m.Content)
	}
	w.WriteHeader(200)
}

func wsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	userid := r.Header.Get("uid")
	uid, _ := strconv.Atoi(userid)
	cl := NewClient(c, uid)
	defer cl.Close()

	for {
		//mt, message, err := cl.c.ReadMessage()
		//if err != nil {
		//	log.Println("read:", err)
		//	break
		//}
		//log.Printf("recv: %s", message)
		message := <-cl.InfoC
		err = c.WriteMessage(1, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/wsocket", wsocket)
	http.HandleFunc("/newmsg", NewMsg)
	Push = NewPushPool()
	log.Println(http.ListenAndServe(":8009", nil))
}
