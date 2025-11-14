package migrations

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		tracer := NewTracer("v0.4.5")
		tracer.Printf("go ...")

		// adapt to new workflow data structure
		{
			type dWorkflowNode struct {
				Id     string           `json:"id"`
				Type   string           `json:"type"`
				Data   map[string]any   `json:"data"`
				Blocks []*dWorkflowNode `json:"blocks,omitempty,omitzero"`
			}

			var deepMigrateNode func(node *dWorkflowNode) (_node *dWorkflowNode, _migrated bool)
			var deepMigrateNodes func(nodes []*dWorkflowNode) (_nodes []*dWorkflowNode, _migrated bool)
			deepMigrateNode = func(node *dWorkflowNode) (*dWorkflowNode, bool) {
				migrated := false

				if node.Type == "bizDeploy" {
					if node.Data != nil {
						if _, ok := node.Data["config"]; ok {
							nodeCfg := node.Data["config"].(map[string]any)
							if nodeCfg["provider"] == "tencentcloud-ssldeploy" {
								if nodeCfg["providerConfig"] != nil {
									providerCfg := nodeCfg["providerConfig"].(map[string]any)
									providerCfg["resourceProduct"] = providerCfg["resourceType"]
									delete(providerCfg, "resourceType")
									nodeCfg["providerConfig"] = providerCfg

									node.Data["config"] = nodeCfg
									migrated = true
								}
							} else if nodeCfg["provider"] == "tencentcloud-sslupdate" {
								if nodeCfg["providerConfig"] != nil {
									providerCfg := nodeCfg["providerConfig"].(map[string]any)
									providerCfg["resourceProducts"] = providerCfg["resourceTypes"]
									delete(providerCfg, "resourceTypes")
									nodeCfg["providerConfig"] = providerCfg

									node.Data["config"] = nodeCfg
									migrated = true
								}
							} else if nodeCfg["provider"] == "aliyun-waf" {
								if nodeCfg["providerConfig"] != nil {
									providerCfg := nodeCfg["providerConfig"].(map[string]any)
									providerCfg["serviceType"] = "cname"
									nodeCfg["providerConfig"] = providerCfg

									node.Data["config"] = nodeCfg
									migrated = true
								}
							}
						}
					}
				}

				if len(node.Blocks) > 0 {
					if newBlocks, changed := deepMigrateNodes(node.Blocks); changed {
						node.Blocks = newBlocks
						migrated = true
					}
				}

				return node, migrated
			}
			deepMigrateNodes = func(nodes []*dWorkflowNode) ([]*dWorkflowNode, bool) {
				migrated := false

				for i, node := range nodes {
					if newNode, changed := deepMigrateNode(node); changed {
						nodes[i] = newNode
						migrated = true
					}
				}

				return nodes, migrated
			}

			// update collection `workflow`
			//   - migrate field `graphDraft` / `graphContent`
			{
				collection, err := app.FindCollectionByNameOrId("tovyif5ax6j62ur")
				if err != nil {
					return err
				}

				records, err := app.FindAllRecords(collection)
				if err != nil {
					return err
				}

				for _, record := range records {
					changed := false

					graphDraft := make(map[string]any)
					if err := record.UnmarshalJSONField("graphDraft", &graphDraft); err == nil {
						if _, ok := graphDraft["nodes"]; ok {
							nodes := make([]*dWorkflowNode, 0)
							if err := mapstructure.Decode(graphDraft["nodes"], &nodes); err != nil {
								return err
							}

							if newNodes, migrated := deepMigrateNodes(nodes); migrated {
								graphDraft["nodes"] = newNodes
								record.Set("graphDraft", graphDraft)
								changed = true
							}
						}
					}

					graphContent := make(map[string]any)
					if err := record.UnmarshalJSONField("graphContent", &graphContent); err == nil {
						if _, ok := graphContent["nodes"]; ok {
							nodes := make([]*dWorkflowNode, 0)
							if err := mapstructure.Decode(graphContent["nodes"], &nodes); err != nil {
								return err
							}

							if newNodes, migrated := deepMigrateNodes(nodes); migrated {
								graphContent["nodes"] = newNodes
								record.Set("graphContent", graphContent)
								changed = true
							}
						}
					}

					if changed {
						if err := app.Save(record); err != nil {
							return err
						}

						tracer.Printf("record #%s in collection '%s' updated", record.Id, collection.Name)
					}
				}

				if _, err := app.DB().NewQuery("UPDATE workflow SET graphDraft = REPLACE(graphDraft, '\"matchPattern\"', '\"domainMatchPattern\"')").Execute(); err != nil {
					return err
				}

				if _, err := app.DB().NewQuery("UPDATE workflow SET graphContent = REPLACE(graphContent, '\"matchPattern\"', '\"domainMatchPattern\"')").Execute(); err != nil {
					return err
				}
			}

			// update collection `workflow_run`
			//   - migrate field `graph`
			{
				collection, err := app.FindCollectionByNameOrId("qjp8lygssgwyqyz")
				if err != nil {
					return err
				}

				records, err := app.FindAllRecords(collection)
				if err != nil {
					return err
				}

				for _, record := range records {
					changed := false

					graph := make(map[string]any)
					if err := record.UnmarshalJSONField("graph", &graph); err == nil {
						if _, ok := graph["nodes"]; ok {
							nodes := make([]*dWorkflowNode, 0)
							if err := mapstructure.Decode(graph["nodes"], &nodes); err != nil {
								return err
							}

							if newNodes, migrated := deepMigrateNodes(nodes); migrated {
								graph["nodes"] = newNodes
								record.Set("graph", graph)
								changed = true
							}
						}
					}

					if changed {
						if err := app.Save(record); err != nil {
							return err
						}

						tracer.Printf("record #%s in collection '%s' updated", record.Id, collection.Name)
					}
				}

				if _, err := app.DB().NewQuery("UPDATE workflow_run SET graph = REPLACE(graph, '\"matchPattern\"', '\"domainMatchPattern\"')").Execute(); err != nil {
					return err
				}
			}

			// update collection `workflow_output`
			//   - migrate field `nodeConfig`
			{
				if _, err := app.DB().NewQuery("UPDATE workflow_output SET nodeConfig = REPLACE(nodeConfig, '\"matchPattern\"', '\"domainMatchPattern\"')").Execute(); err != nil {
					return err
				}

				if _, err := app.DB().NewQuery("UPDATE workflow_output SET nodeConfig = REPLACE(nodeConfig, '\"resourceType\"', '\"resourceProduct\"') WHERE nodeConfig LIKE '%\"provider\":\"tencentcloud-ssldeploy\"%' OR nodeConfig LIKE '%\"provider\":\"tencentcloud-sslupdate\"%'").Execute(); err != nil {
					return err
				}
			}
		}

		tracer.Printf("done")
		return nil
	}, func(app core.App) error {
		return nil
	})
}
