package bleservice

import "fmt"

func strToSlice(str string, dst []byte) error {
	capacity := len(dst)
	if len(str) > capacity {
		return fmt.Errorf("str too long")
	}

	dst[0] = uint8(len(str))
	for i := 0; i < len(str); i++ {
		dst[i+1] = str[i]
	}

	return nil
}

func strTo33Bytes(str string) ([33]byte, error) {
	dst := [33]byte{}
	return dst, strToSlice(str, dst[:])
}

func strTo65Bytes(str string) ([65]byte, error) {
	dst := [65]byte{}
	return dst, strToSlice(str, dst[:])
}

func strTo129Bytes(str string) ([129]byte, error) {
	dst := [129]byte{}
	return dst, strToSlice(str, dst[:])
}

func bytesToStr(strBytes []byte) (string, error) {
	l := int(strBytes[0])
	if l > len(strBytes)-1 {
		return "", fmt.Errorf("length not correctly encoded")
	}
	return string(strBytes[1 : l+1]), nil

}
