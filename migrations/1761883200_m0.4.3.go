package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		tracer := NewTracer("v0.4.3")
		tracer.Printf("go ...")

		// update collection `certificate`
		//   - rename field `acmeRenewed` to `isRenewed`
		//   - add field `isRevoked`
		{
			collection, err := app.FindCollectionByNameOrId("4szxr9x43tpj6np")
			if err != nil {
				return err
			}

			if err := collection.Fields.AddMarshaledJSONAt(14, []byte(`{
				"hidden": false,
				"id": "bool810050391",
				"name": "isRenewed",
				"presentable": false,
				"required": false,
				"system": false,
				"type": "bool"
			}`)); err != nil {
				return err
			}

			if err := collection.Fields.AddMarshaledJSONAt(15, []byte(`{
				"hidden": false,
				"id": "bool3680845581",
				"name": "isRevoked",
				"presentable": false,
				"required": false,
				"system": false,
				"type": "bool"
			}`)); err != nil {
				return err
			}

			if err := app.Save(collection); err != nil {
				return err
			}
		}

		tracer.Printf("done")
		return nil
	}, func(app core.App) error {
		return nil
	})
}
