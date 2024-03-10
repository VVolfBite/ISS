
package util

import (
	"bytes"
	"encoding/gob"
	// "fmt"
	// "github.com/kelindar/bitmap"
	// "math/rand"
	// "net"
	// "reflect"
	// "strconv"
	// "time"
)

func SizeOf(v interface{}) int {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(v)
	return buffer.Len()
}