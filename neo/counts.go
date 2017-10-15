package neo

import (
	"fmt"
)

type nodesCounter struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

func (db *DB) Counts() ([]*nodesCounter, error) {
	rows, err := db.Query(statement)
	if err != nil {
		return nil, fmt.Errorf("read query failed: %v", err)
	}
	defer rows.Close()

	counts := make([]*nodesCounter, 0)
	for rows.Next() {
		item := new(nodesCounter)
		err = rows.Scan(&item.Name, &item.Count)
		if err != nil {
			return nil, fmt.Errorf("failed when iterating rows: %v", err)
		}
		counts = append(counts, item)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error when out from rows: %v", err)
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
