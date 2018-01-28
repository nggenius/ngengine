package share

type ServiceId uint16

const (
	SID_MAX        = 0xFFFF // service id max
	MB_FLAG_APP    = iota   // app
	MB_FLAG_CLIENT          // client
)
