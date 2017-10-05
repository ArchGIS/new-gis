package routes

import "github.com/ArchGIS/new-gis/neo"

type (
	request struct {
		Lang string `query:"lang"`
	}

	Env struct {
		neo.DataStore
	}
)

var db *Env

func InitEnv(source string) error {
	dbInstance, err := neo.InitDB(source)
	if err != nil {
		return err
	}

	db = &Env{dbInstance}
	return nil
}
