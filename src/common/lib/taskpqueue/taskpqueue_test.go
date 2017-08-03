package taskpqueue

import (
	"fmt"
	"testing"
	"time"
)

func Test_taskqueue(b *testing.T) {
	jobtask := func(interval interface{}) interface{} {
		fmt.Println("come on baby go go go !!!", interval)
		duration, err := time.ParseDuration(fmt.Sprintf("%vs", interval))
		if err != nil {
			return fmt.Errorf("time.ParseDuration err :", err)
		}
		time.Sleep(duration)
		return nil
	}
	tq := NewTaskQueue(30, 10, 1000*60*4, jobtask)
	tq.Start()
	for i := 1; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		err := tq.AddTask(i, key, i)
		if err != nil {
			fmt.Println("sorry baby you cant !!")
		}
		fmt.Println("baby baby come in!!!", i)
	}
	time.Sleep(500 * time.Second)
}
