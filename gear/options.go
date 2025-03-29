package gear

// TODO: better handling of index
func GetOptionString(index int, defaultValue string, options ...string) string {
	if len(options)==0 { 
		return defaultValue
		
	} else {
		if index > len(options) {
			return defaultValue
		}
	}
	return options[index]
}
