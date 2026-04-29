package engine

func recoverToValue(input any, result *any) {
	if recover() != nil {
		*result = fallbackValue(input)
	}
}

func recoverToJSONFallback(result *[]byte) {
	if recover() != nil {
		*result = cloneJSONFallback()
	}
}

func recoverToYAMLFallback(result *[]byte) {
	if recover() != nil {
		*result = cloneYAMLFallback()
	}
}
