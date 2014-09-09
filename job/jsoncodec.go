package job

import (
	"encoding/json"
	"io"
	"os"
)

// NewJSONEncoder returns a new encoder that writes to w
func NewJSONEncoder(w io.Writer) *json.Encoder {
	if w == nil {
		w = os.Stdout
	}
	return json.NewEncoder(w)
}
