//go:build go1.21 && !without_contextjson

package json

var (
	Marshal    = json.Marshal
	Unmarshal  = json.Unmarshal
	NewEncoder = json.NewEncoder
	NewDecoder = json.NewDecoder
)

type (
	Encoder     = json.Encoder
	Decoder     = json.Decoder
	Token       = json.Token
	Delim       = json.Delim
	SyntaxError = json.SyntaxError
	RawMessage  = json.RawMessage
)
