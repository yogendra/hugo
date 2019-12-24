package transform

import (
	"bytes"
	"compress/zlib"
	"github.com/spf13/cast"
	"strings"
)

// PlantumlTextFormat renders PlantUML stuff
func (ns *Namespace) PlantumlTextFormat(s interface{}) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}
	ss = strings.TrimSpace(ss)
	compressed, err := deflate([]byte(ss))
	if err != nil {
		return "", err
	}
	return base64Encode(compressed), nil
}

func deflate(input []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := zlib.NewWriterLevel(&b, zlib.BestCompression)
	if err != nil {
		return nil, err
	}
	w.Write(input)
	w.Close()
	return b.Bytes(), nil
}

func base64Encode(input []byte) string {
	var buffer bytes.Buffer
	inputLength := len(input)
	for i := 0; i < 3-inputLength%3; i++ {
		input = append(input, byte(0))
	}

	for i := 0; i < inputLength; i += 3 {
		b1, b2, b3, b4 := input[i], input[i+1], input[i+2], byte(0)

		b4 = b3 & 0x3f
		b3 = ((b2 & 0xf) << 2) | (b3 >> 6)
		b2 = ((b1 & 0x3) << 4) | (b2 >> 4)
		b1 = b1 >> 2

		for _, b := range []byte{b1, b2, b3, b4} {
			buffer.WriteByte(byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"[b]))
		}
	}
	return string(buffer.Bytes())
}
