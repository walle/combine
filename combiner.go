package combine

// Combiner is responsible for initializing, processing and writing a Includer.
// See StreamCombiner for a basic implementation.
type Combiner interface {
	Read(Includer) (Includer, error)
	Combine(Includer) []error
	SetResult([]byte)
	Result() []byte
	Write() error
}
