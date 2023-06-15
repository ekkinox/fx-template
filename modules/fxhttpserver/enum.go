package fxhttpserver

type MiddlewareKind int

const (
	GlobalUse MiddlewareKind = iota
	GlobalPre
	Attached
)

func (k MiddlewareKind) String() string {
	switch k {
	case GlobalUse:
		return "global-use"
	case GlobalPre:
		return "global-pre"
	case Attached:
		return "attached"
	default:
		return "global-use"
	}
}
