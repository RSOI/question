package view

func fieldsToString(fields []string) string {
	res := ""
	for index, r := range fields {
		if index > 0 {
			res += ", "
		}
		res += r
	}
	return res
}
