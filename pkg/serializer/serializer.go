package serializer

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Serialize[K any](obj K) ([]byte, error) {
	return json.Marshal(obj)
}

func Deserialize[K any](b []byte) (obj K, err error) {
	err = json.Unmarshal(b, &obj)
	if err != nil {
		return obj, err
	}
	return
}

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func NewDecoder(r io.Reader) *jsoniter.Decoder {
	return json.NewDecoder(r)
}

func NewEncoder(w io.Writer) *jsoniter.Encoder {
	return json.NewEncoder(w)
}
