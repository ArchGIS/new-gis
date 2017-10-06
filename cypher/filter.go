package cypher

import "fmt"

// Filter fills statement with filter
func Filter(statement, filter string) string {
	return fmt.Sprintf(statement, filter)
}
