package ossfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateImage(t *testing.T) {
	name := "logo.png"
	fin, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	defer fin.Close()

	data, err := ioutil.ReadAll(fin)
	if err != nil {
		t.Fatal(err)
	}

	url, err := PutFile("test", name, data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(url)
}
