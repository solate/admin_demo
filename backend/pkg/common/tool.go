package common

import (
	"math/rand"
	"time"
)

func StrDateToTimestamp(fromDate string) (int64, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, fromDate)
	if err != nil {
		return 0, err
	}
	unixTimestamp := t.UnixMilli()
	return unixTimestamp, nil
}

// 随机生成6位数
func GenerateSixDigitNumber() int {
	return rand.Intn(900000) + 100000
}

// 随机生成字符串
func GenRandomString(n int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

//
