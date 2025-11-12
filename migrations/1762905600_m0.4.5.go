package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		tracer := NewTracer("v0.4.5")
		tracer.Printf("go ...")

		// adapt to new workflow data structure
		{
			if _, err := app.DB().NewQuery("UPDATE workflow SET graphDraft = REPLACE(graphDraft, '\"matchPattern\"', '\"domainMatchPattern\"')").Execute(); err != nil {
				return err
			}

			if _, err := app.DB().NewQuery("UPDATE workflow SET graphContent = REPLACE(graphContent, '\"matchPattern\"', '\"domainMatchPattern\"')").Execute(); err != nil {
				return err
			}

			if _, err := app.DB().NewQuery("UPDATE workflow_run SET graph = REPLACE(graph, '\"matchPattern\"', '\"domainMatchPattern\"')").Execute(); err != nil {
				return err
			}

			if _, err := app.DB().NewQuery("UPDATE workflow_output SET nodeConfig = REPLACE(nodeConfig, '\"matchPattern\"', '\"domainMatchPattern\"')").Execute(); err != nil {
				return err
			}
		}

		tracer.Printf("done")
		return nil
	}, func(app core.App) error {
		return nil
	})
}
