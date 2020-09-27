package utility

// Contains checks a string slice for membership
func Contains(ss []string, s string) bool {
	for _, m := range ss {
		if m == s {
			return true
		}
	}
	return false
}
