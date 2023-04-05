package utils

func Map[E any, F any](entries []E, fn func(E) F) []F {
	var res = make([]F, len(entries))
	for i, entry := range entries {
		res[i] = fn(entry)
	}

	return res
}

func MapWithErr[E any, F any](entries []E, fn func(E) (F, error)) ([]F, error) {
	var res = make([]F, len(entries))
	var err error
	for i, entry := range entries {
		res[i], err = fn(entry)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}