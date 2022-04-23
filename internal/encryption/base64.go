package encryption

import (
	"log"
)

const padchar = '='
const alpha = "LVoJPiCN2R8G90yg+hmFHuacZ1OWMnrsSTXkYpUq/3dlbfKwv6xztjI7DeBE45QA"

func getBytes(s string, i int) int {
	b := s[i]
	if b > 255 {
		log.Fatalln("INVALID_CHARACTER_ERR: DOM Exception 5")
	}
	return int(b)
}

func GetBase64(s string) string {
	x := make([]byte, 0)
	imax := len(s) - len(s)%3
	if len(s) == 0 {
		return s
	}

	for i := 0; i < imax; i += 3 {
		b10 := (getBytes(s, i) << 16) | (getBytes(s, i+1) << 8) | getBytes(s, i+2)
		x = append(x, alpha[(b10>>18)], alpha[((b10>>12)&63)], alpha[((b10>>6)&63)], alpha[(b10&63)])
	}

	i := imax
	if len(s)-imax == 1 {
		b10 := getBytes(s, i) << 16
		x = append(x, alpha[(b10>>18)], alpha[((b10>>12)&63)], padchar, padchar)
	} else if len(s)-imax == 2 {
		b10 := (getBytes(s, i) << 16) | (getBytes(s, i+1) << 8)
		x = append(x, alpha[(b10>>18)], alpha[((b10>>12)&63)], alpha[((b10>>6)&63)], padchar)
	}

	return string(x)
}
