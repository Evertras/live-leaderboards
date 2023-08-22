package server

func ptr[K any](item K) *K {
	return &item
}
