package utils

import (
	"math/rand"
	"time"
)

func GetUuid(l int) string {
	initByte := []byte{}
	result := []byte{}

	for i := 48; i < 123; i++ {
		switch {
		case i < 58:
			initByte = append(initByte, byte(i))
		case i >= 65 && i < 91:
			initByte = append(initByte, byte(i))
		case i >= 97 && i < 123:
			initByte = append(initByte, byte(i))
		}

	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < l; i++ {
		result = append(result, initByte[rand.Intn(len(initByte))])
	}
	return string(result)
}

func GetLowerAndNumUuid(l int) string {
	initByte := []byte{}
	result := []byte{}

	for i := 48; i < 123; i++ {
		switch {
		case i < 58:
			initByte = append(initByte, byte(i))
		case i >= 97 && i < 123:
			initByte = append(initByte, byte(i))
		}

	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < l; i++ {
		result = append(result, initByte[rand.Intn(len(initByte))])
	}
	return string(result)
}
