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
	WithSourceCodeInfo      bool
	WithJsonTag             bool
	WithGoogleProtobuf      bool
}

func LoadOptions(ops ...OptionFunc) *Option {
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

func WithSourceCodeInfo() OptionFunc {
	return func(option *Option) {
		option.WithSourceCodeInfo = true
	}
}

func WithJsonTag() OptionFunc {
	return func(option *Option) {
		option.WithJsonTag = true
	}
}

func WithMessageType(t MessageType) OptionFunc {
	return func(option *Option) {
		option.MessageType = t
	}
}

func WithGoogleProtobuf() OptionFunc {
	return func(option *Option) {
		option.WithGoogleProtobuf = true
	}
}

type IDLConfig struct {
	Main        string
	IDLs        map[string][]byte
	IncludePath []string
}
