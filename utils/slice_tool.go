package utils

func Contains(strarr []string, str string) bool {
	if len(strarr) == 0 {
		return false
	}
	for _, s := range strarr {
		if s == str {
			return true
		}
	}
	return false
}
