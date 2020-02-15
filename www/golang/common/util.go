package common

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func Contains(a []string, x string) (bool, int) {
	for key, n := range a {
		if x == n || strings.ToUpper(x) == n {
			return true, key
		}
	}
	return false, -1
}

func CallAPI(url string, method string, body []byte, header map[string]string) {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	for key, value := range header {
		req.Header.Set(key, value)
	}

	clt := http.Client{}
	resp, respErr := clt.Do(req)
	if respErr != nil {
		Log.Error("call response api error!")
	}
	fmt.Println(resp)
	defer resp.Body.Close()
}

func CheckGoRoutineNum() {
	fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine())
}

func Mapkey(m map[string][]string, value string) (key string, ok bool) {
	for k, v := range m {
		for _, sliceV := range v {
			if sliceV == value {
				key = k
				ok = true
				return
			}
		}
	}
	return
}

func StringToTime(strTime string) time.Time {
	if strings.Contains(strTime, "/") {
		layout := "2006/01/02 15:04"
		formatLayout := "2006-01-02 15:04:05"
		layoutTime, _ := time.Parse(layout, strTime)
		formatLayoutTimeStr := layoutTime.Format(formatLayout)
		formatLayoutTime, _ := time.Parse(formatLayout, formatLayoutTimeStr)
		return formatLayoutTime
	} else {
		formatLayout := "2006-01-02 15:04"
		layoutTime, _ := time.Parse(formatLayout, strTime)
		return layoutTime
	}
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
