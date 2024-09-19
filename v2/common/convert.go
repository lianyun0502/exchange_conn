package common


import (
	"unsafe"
)
func String2Bytes(s string) []byte {
    if s == "" {
		return nil
	}
    return unsafe.Slice(unsafe.StringData(s), len(s))
}

func Bytes2String(b []byte) string {
	if len(b) == 0 {
		return	""
	}
    return unsafe.String(unsafe.SliceData(b), len(b))
}