package neo

type epochProps struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// func (db *DB) Epochs(req gin.H) ([]epochProps, error) {
// 	rows, err := db.QueryNeo(
// 		epochsStatement,
// 		gin.H{"language": req["lang"]},
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	data, _, err := rows.All()
// 	if err != nil {
// 		return nil, err
// 	}

// 	epochs := make([]epochProps, len(data))
// 	for i, row := range data {
// 		epochs[i] = epochProps{
// 			ID:   row[0].(int64),
// 			Name: row[1].(string),
// 		}
// 	}

// 	return epochs, nil
// }

const (
	epochsStatement = `
		MATCH (n:Epoch)-[:translation {lang: {language}}]->(tr:Translate)
		RETURN n.id as id, tr.name as name
	`
)
