package internal

func Values(m map[int64]string) []string {
	values := make([]string, len(m))
	i := 0
	for _, v := range m {
		values[i] = v
		i++
	}

	return values
}
