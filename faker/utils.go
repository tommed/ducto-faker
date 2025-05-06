package faker

import "fmt"

func QuoteString(val string) string {
	return fmt.Sprintf(`"%s"`, val)
}
