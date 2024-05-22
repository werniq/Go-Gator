package utils

// Split splits strings with separator sep
// because strings.Split returns a slice with length 1 if s is empty string,
// I decided to write my own split function which will return empty array in that case
func Split(s, sep string) []string {
	var result []string

	if len(s) == 0 {
		return result
	}

	sepLen := len(sep)
	if sepLen == 0 {
		for _, ch := range s {
			result = append(result, string(ch))
		}
		return result
	}

	start := 0
	for i := 0; i+sepLen <= len(s); {
		if s[i:i+sepLen] == sep {
			result = append(result, s[start:i])
			i += sepLen
			start = i
		} else {
			i++
		}
	}

	result = append(result, s[start:])

	return result
}
