package main

import (
	"encoding/json"
	"fmt"
	"github.com/sishen007/gocrontab/common"
)

func main() {
	var (
		oldJobVal common.Job
		err       error
	)
	str := `{"name":"job1","command":"echo hello","cronExpr":"* * * * *"}`
	if err = json.Unmarshal([]byte(str), &oldJobVal); err != nil {
		return
	}
	fmt.Println("值:", oldJobVal)
}
