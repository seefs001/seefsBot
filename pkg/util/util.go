package util

import "math/rand"

// GenRandomString 生成随机字符串
// length 生成长度
// specialChar 是否生成特殊字符
func GenRandomString(length int, specialChar bool) string {

	letterBytes := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	special := "!@#%$*.="

	if specialChar {
		letterBytes = letterBytes + special
	}

	chars := []byte(letterBytes)

	if length == 0 {
		return ""
	}

	clen := len(chars)
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			return ""
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue // Skip this number to avoid modulo bias.
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
