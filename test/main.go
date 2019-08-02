package main

import (
	"fmt"
	"time"
)

type T struct {
	Name string
	Age  int64
}

func main() {
	//var (
	//	oldJobVal common.Job
	//	err       error
	//)
	//str := `{"name":"job1","command":"echo hello","cronExpr":"* * * * *"}`
	//if err = json.Unmarshal([]byte(str), &oldJobVal); err != nil {
	//	return
	//}
	//fmt.Println("值:", oldJobVal)

	// test Time
	//t := time.Now().Add(10 * time.Second).Sub(time.Now())
	//fmt.Printf("%#v", t)

	t := &T{"张三", 20}
	i := 1
	modifyT(t, &i)
	go getT(t, &i)
	time.Sleep(2 * time.Second)
	fmt.Println("t:", t)
	fmt.Println("i:", i)
}
func modifyT(t *T, i *int) {
	go func() {
		t.Name = "李四"
		*i = 2
	}()
}
func getT(t *T, i *int) {
	fmt.Println("getT t:", t)
	fmt.Println("getT i:", *i)
}
