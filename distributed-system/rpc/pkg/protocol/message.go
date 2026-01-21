package protocol

type Request struct {
	ID            uint64
	ServiceMethod string
	Args          []byte
}

type Response struct {
	ID   uint64
	Data []byte
	Err  string
}
