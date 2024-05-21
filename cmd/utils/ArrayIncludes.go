package utils

func ArrayIncludes(arr []string, s string) bool {
	if arr == nil {
		return true
	}

	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}
