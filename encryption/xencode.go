package encryption

import (
	"math"
	"strings"
)

func force(msg string) []byte {
	return []byte(msg)
}

func ordat(msg string, idx int) int {
	if len(msg) > idx {
		return int(msg[idx])
	}
	return 0
}

func senCode(msg string, key bool) []int {
	l := len(msg)
	pwd := make([]int, 0)
	for i := 0; i < l; i += 4 {
		pwd = append(pwd, ordat(msg, i)|ordat(msg, i+1)<<8|ordat(msg, i+2)<<16|ordat(msg, i+3)<<24)
	}
	if key {
		pwd = append(pwd, l)
	}

	return pwd
}

func lenCode(msg []int) string {
	l := len(msg)
	res := make([]string, l)
	for i := 0; i < l; i++ {
		res[i] = string([]byte{byte(msg[i] & 0xff), byte(msg[i] >> 8 & 0xff), byte(msg[i] >> 16 & 0xff), byte(msg[i] >> 24 & 0xff)})
	}

	return strings.Join(res, "")
}

func GetXencode(msg, key string) string {
	if msg == "" {
		return ""
	}
	pwd := senCode(msg, true)
	pwdk := senCode(key, false)
	if len(pwdk) < 4 {
		for i := 0; i < 4-len(pwdk); i++ {
			pwdk = append(pwdk, 0)
		}
	}
	n := len(pwd) - 1
	z := pwd[n]
	y := pwd[0]
	c := 0x86014019 | 0x183639A0
	m := 0
	e := 0
	p := 0
	q := int(math.Floor(6 + 52/(float64(n)+1)))
	d := 0
	for i := q; i > 0; i-- {
		d = (d + c) & (0x8CE0D9BF | 0x731F2640)
		e = d >> 2 & 3
		p = 0
		for p < n {
			y = pwd[p+1]
			m = z>>5 ^ y<<2
			m = m + ((y>>3 ^ z<<4) ^ (d ^ y))
			m = m + (pwdk[(p&3)^e] ^ z)
			pwd[p] = (pwd[p] + m) & (0xEFB8D130 | 0x10472ECF)
			z = pwd[p]
			p = p + 1
		}
		y = pwd[0]
		m = z>>5 ^ y<<2
		m = m + ((y>>3 ^ z<<4) ^ (d ^ y))
		m = m + (pwdk[(p&3)^e] ^ z)
		pwd[n] = (pwd[n] + m) & (0xBB390742 | 0x44C6F8BD)
		z = pwd[n]
	}

	return lenCode(pwd)
}
