package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"reflect"
    "time"

	"github.com/google/uuid"
)


func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

func GetUUID() string  {
    return uuid.NewString()
}

func GetSignature(secretKey string, data string) string {
    mac := hmac.New(sha256.New, []byte(secretKey))
    mac.Write([]byte(data))
    return hex.EncodeToString(mac.Sum(nil))
} 

func IsEmpty(v interface{}) bool {
    // 使用 reflect 處理不同類型的「空值」判斷
    switch reflect.ValueOf(v).Kind() {
    case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
        // 對於指針、切片、映射、通道、函數和接口，使用 IsNil 判斷是否為 nil
        return reflect.ValueOf(v).IsNil()
    case reflect.String:
        // 對於字符串，判斷是否為空字符串
        return v == ""
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        // 對於整數類型，判斷是否為 0
        return v == 0
    case reflect.Float32, reflect.Float64:
        // 對於浮點數類型，判斷是否為 0.0
        return v == 0.0
    default:
        // 對於其他類型，返回 false 表示無法判斷或不為「空值」
        return false
    }
}

type Timer struct {
	Interval time.Duration
	Handler func()

	stop chan struct{}
	reset chan struct{}
}
func (hbt *Timer) Start(handle func()) {
	if handle != nil {
		hbt.Handler = handle
	}
	hbt.stop = make(chan struct{})
	hbt.reset = make(chan struct{})
	go func() {
		for {
			select {
            case <-time.After(hbt.Interval):
                hbt.Handler()
            case <-hbt.stop:
                return
			case <-hbt.reset:
				continue
            }
		}
	}()
}
func (hbt *Timer) Stop() {
	hbt.stop <- struct{}{}
}	

func (hbt *Timer) Reset() {
	hbt.reset <- struct{}{}
}