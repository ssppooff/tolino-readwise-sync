package utils

func Map[E any, F any](sl []E, fn func(E) F) []F {
	var res = make([]F, len(sl))
	for i, el := range sl {
		res[i] = fn(el)
	}

	return res
}

func MapWithErr[E any, O any](sl []E, fn func(E) (O, error)) ([]O, error) {
	var res = make([]O, len(sl))
	var err error
	for i, el := range sl {
		res[i], err = fn(el)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func Filter[T any](sl []T, p func(T) bool) (res []T) {
	for _, el := range sl {
		if p(el) {
			res = append(res, el)
		}
	}

	return
}

func Not[T any](fn func(T) bool) func(T) bool {
	return func(e T) bool {
		return !fn(e)
	}
}