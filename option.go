package protobuf

type OptionFunc func(*Option)

type MessageType int

const (
	MessageType_PB = iota
	MessageType_JSON
)

type Option struct {
	MessageType             MessageType
	RequireSyntaxIdentifier bool
}

func loadOptions(ops ...OptionFunc) *Option {
	rr := new(Option)
	for _, elem := range ops {
		elem(rr)
	}
	return rr
}

func WithRequireSyntaxIdentifier() OptionFunc {
	return func(option *Option) {
		option.RequireSyntaxIdentifier = true
	}
}

func WithMessageType(t MessageType) OptionFunc {
	return func(option *Option) {
		option.MessageType = t
	}
}
