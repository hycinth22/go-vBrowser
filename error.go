package vBrowser

type Op string

const (
	ParseURL      Op = "parseURL"
	CreateRequest Op = "createRequest"
	DoRequest     Op = "doRequest"
	// ParseDocument Op = "parseDocument"
)

type Error struct {
	Op  Op
	Err error
}

func (e *Error) Error() string {
	return "Failed to execute the operation " + string(e.Op)
}
