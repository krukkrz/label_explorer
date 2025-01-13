package handlers

func validateOrder(value string) bool {
	allowedValues := map[string]struct{}{
		"asc":  {},
		"desc": {},
	}

	_, ok := allowedValues[value]
	return ok
}
