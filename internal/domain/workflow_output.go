package domain

const CollectionNameWorkflowOutput = "workflow_output"

type WorkflowOutput struct {
	Meta
	WorkflowId string           `json:"workflowId" db:"workflowRef"`
	RunId      string           `json:"runId" db:"runRef"`
	NodeId     string           `json:"nodeId" db:"nodeId"`
	Node       *WorkflowNode    `json:"node" db:"node"`
	Outputs    []WorkflowNodeIO `json:"outputs" db:"outputs"` // Deprecated
	Succeeded  bool             `json:"succeeded" db:"succeeded"`
}
