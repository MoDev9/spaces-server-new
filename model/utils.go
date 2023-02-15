package model

import (
	"bytes"
	"encoding/base32"
	"encoding/json"
	"io"
	"strconv"
	"time"
	"unicode"

	"github.com/pborman/uuid"
)

const (
	MINIMUM_LENGTH_PASSWORD = 8

	HEADER_REQUESTED_WITH     = "X-Requested-With"
	HEADER_REQUESTED_WITH_XML = "XMLHttpRequest"
	HEADER_CSRF_TOKEN         = "X-CSRF-Token"
)

//Get milliseconds since epoch
func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func IdToString(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

func IsValidId(id string) bool {
	if len(id) != 26 {
		return false
	}

	for _, v := range id {
		if !unicode.IsLetter(v) && !unicode.IsNumber(v) {
			return false
		}
	}

	return true
}

func NewAppError() {

}

func JsonToMap(r io.Reader) map[string]interface{} {
	var objmap map[string]interface{}
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(&objmap); err != nil {
		return make(map[string]interface{})
	}
	return objmap
}

func MapToJson(obj map[string]string) string {
	b, _ := json.Marshal(obj)
	return string(b)
}

func StringInterfaceToJson(objmap map[string]interface{}) string {
	b, _ := json.Marshal(objmap)
	return string(b)
}

func ArrayToJson(objmap []string) string {
	b, _ := json.Marshal(objmap)
	return string(b)
}

func ArrayFromJson(data io.Reader) []string {
	decoder := json.NewDecoder(data)

	var arr []string
	if err := decoder.Decode(&arr); err != nil {
		return make([]string, 0)
	}
	return arr
}

func NewId() string {
	var b bytes.Buffer
	encoder := base32.NewEncoder(base32.StdEncoding, &b)
	encoder.Write(uuid.NewRandom())
	encoder.Close()
	b.Truncate(26) // removes the '==' padding
	return b.String()
}
