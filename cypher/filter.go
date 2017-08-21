package cypher

import "fmt"

// Filter fills statement with filter
func Filter(statement, filter string) string {
	return fmt.Sprintf(statement, filter)
}

// BuildRegexpFilter return neo4j regexp filter
// for case-insensitive text search
func BuildRegexpFilter(needle interface{}) string {
	return fmt.Sprintf("(?ui).*%s.*$", needle)
}
