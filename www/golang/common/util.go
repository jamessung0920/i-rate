package common

import (
	"fmt"
	"bytes"
	"strings"
	"net/http"
	"runtime"
)

func Contains(a []string, x string) bool {
    for _, n := range a {
        if x == n || strings.ToUpper(x) == n {
            return true
        }
    }
    return false
}

func CallAPI(url string, method string, body []byte) {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	clt := http.Client{}
	resp, respErr := clt.Do(req)
	if respErr != nil {
		panic(respErr)
		fmt.Println(respErr)
	}
	defer resp.Body.Close()
}

func CheckGoRoutineNum() {
	fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine())
}