package routes

import "github.com/ArchGIS/new-gis/neo"

type (
	request struct {
		Lang string `query:"lang"`
	}

	Env struct {
		db neo.DataStore
	}
)

var Model *Env

func InitEnv(source string) error {
	db, err := neo.InitDB(source)
	if err != nil {
		return err
	}

	Model = &Env{db}
	return nil
}
