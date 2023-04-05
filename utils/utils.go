package utils

func Map[E any, F any](entries []E, fn func(E) F) []F {
	var res = make([]F, len(entries))
	for i, entry := range entries {
		res[i] = fn(entry)
	}
	return res
}
