package types

type Event struct {
	ID     string
	Stream string
	Blob   []byte
	Meta   map[string]string
}
