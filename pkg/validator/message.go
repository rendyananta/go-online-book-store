package validator

import "fmt"

func ErrMessage(errTag string, fieldName string) string {
	switch errTag {
	case "required":
		return fmt.Sprintf("field [%s] is required", fieldName)
	case "email":
		return fmt.Sprintf("field [%s] must be a valid email", fieldName)
	default:
		return fmt.Sprintf("field [%s] is invalid", fieldName)
	}
}
