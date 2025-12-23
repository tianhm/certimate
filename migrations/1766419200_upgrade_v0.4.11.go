package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// clean old migrations
		{
			migrations := []string{
				"1757476800_m0.4.0_migrate.go",
				"1757476801_m0.4.0_initialize.go",
				"1760486400_m0.4.1.go",
				"1762142400_m0.4.3.go",
				"1762516800_m0.4.4.go",
				"1763373600_m0.4.5.go",
				"1763640000_m0.4.6.go",
			}
			for _, name := range migrations {
				app.DB().NewQuery("DELETE FROM _migrations WHERE file='" + name + "'").Execute()
			}
		}

		return nil
	}, func(app core.App) error {
		return nil
	})
}
