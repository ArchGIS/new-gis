package cypherBuilder

import (
	"fmt"
)

const (
	statement = `
		OPTIONAL MATCH (s:%s {id: %d})-[:has]->(sr:SpatialReference)-[:has]->(srt:SpatialReferenceType)
		WITH sr, srt%s
		ORDER BY srt.id ASC, sr.date DESC LIMIT 1
		WITH %scollect({x: sr.x, y: sr.y}) AS rows`

	ending = `
		UNWIND rows AS row
		RETURN row.x as x, row.y as y`
)

const (
	// First statement or not
	firstStmt = iota
	otherStmt
)

func BuildCoordinates(ids []int, entity string) string {
	var result string
	counter := 0

	for _, stmt := range ids {
		switch counter {
		case firstStmt:
			result += fmt.Sprintf(statement, entity, stmt, "", "")
		default:
			result += fmt.Sprintf(statement, entity, stmt, ", rows", "rows + ")
		}
		counter++
	}
	result += ending

	return result
}

func removeDuplicates(a []int) []int {
	result := []int{}
	seen := map[int]int{}
	for _, val := range a {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}
