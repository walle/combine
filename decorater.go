package combine

// Decorater is an interface that supplies the strucure to extend/modify output from combiners.
type Decorater interface {
	Decorate(Combiner) error
}
