package cypherBuilder

// func TestBuild(t *testing.T) {
// 	ids := []int{1, 4, 1000}
// 	entity := "Monument"

// 	expected := `
// 		MATCH (s:Monument {id: 1})-[:has]->(sr:SpatialReference)-[:has]->(srt:SpatialReferenceType)
// 		WITH s, sr, srt
// 		ORDER BY srt.id ASC, sr.date DESC LIMIT 1
// 		WITH collect({date: sr.date, x: sr.x, y: sr.y, type: srt.name}) AS rows
// 		MATCH (s:Monument {id: 4})-[:has]->(sr:SpatialReference)-[:has]->(srt:SpatialReferenceType)
// 		WITH s, sr, srt, rows
// 		ORDER BY srt.id ASC, sr.date DESC LIMIT 1
// 		WITH rows + collect({date: sr.date, x: sr.x, y: sr.y, type: srt.name}) AS rows
// 		MATCH (s:Monument {id: 1000})-[:has]->(sr:SpatialReference)-[:has]->(srt:SpatialReferenceType)
// 		WITH s, sr, srt, rows
// 		ORDER BY srt.id ASC, sr.date DESC LIMIT 1
// 		WITH rows + collect({date: sr.date, x: sr.x, y: sr.y, type: srt.name}) AS rows
// 		UNWIND rows AS row
// 		RETURN row.date as date, row.x as x, row.y as y, row.type as type`

// 	if result := BuildCoordinates(ids, entity); result != expected {
// 		t.Errorf("\nReceived:%s\n != \nExpected:%s", result, expected)
// 	}
// }
