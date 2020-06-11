package errorcode

type ErrorCode int

const (
	Success ErrorCode = iota
	ConnectionClosed
	ConnectionTimeout
	ConnectionUnknownError
)

func (ec ErrorCode) string() string {
	str := ""
	switch ec {
	case Success:
		str = "Success"
	case ConnectionClosed:
		str = "ConnectionClosed"
	case ConnectionTimeout:
		str = "ConnectionTimeout"
	case ConnectionUnknownError:
		str = "ConnectionUnknownError"
	}
	return str
}