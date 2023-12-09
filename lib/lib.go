package lib

func Must[T any](val T, err any) T {
	if err != nil {
		panic(err)
	}

	return val
}
