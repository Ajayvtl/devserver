package legacy

func NewNoop(logger ...interface{}) Runtime {
	return &localAdapter{} // Reusing adapter logic since it's a dummy
}
