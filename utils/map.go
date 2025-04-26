package utils

func Map[K, V any](original []K, mapperFn func(K) V) []V {
	res := make([]V, len(original))
	for i, it := range original {
		res[i] = mapperFn(it)
	}
	return res
}
