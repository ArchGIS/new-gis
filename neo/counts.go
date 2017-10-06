package neo

import (
	"github.com/gin-gonic/gin"
)

type NodesCounter struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

func (db *DB) Counts() ([]NodesCounter, error) {
	rows, err := db.QueryNeo(statement, gin.H{})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, _, err := rows.All()
	if err != nil {
		return nil, err
	}

	counts := make([]NodesCounter, len(data))
	for i, row := range data {
		counts[i] = NodesCounter{
			Name:  row[0].(string),
			Count: row[1].(int64),
		}
	}

	return counts, nil
}

const (
	statement = `
		OPTIONAL MATCH (a:Author)
		RETURN labels(a)[0] as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Research)
		RETURN labels(a)[0] as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Heritage)
		RETURN labels(a)[0] as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Monument)
		RETURN labels(a)[0] as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Artifact)
		RETURN labels(a)[0] as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Radiocarbon)
		RETURN labels(a)[0] as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Report)
		RETURN labels(a)[0] as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Monography)
		RETURN "Monography" as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Article)
		RETURN "Article" as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:ArchiveDoc)
		RETURN "ArchiveDoc" as name, count(a) as count
	`
)
