package engine

import (
	"github.com/certimate-go/certimate/internal/domain"
)

type Node = domain.WorkflowNode

type NodeType = domain.WorkflowNodeType

const (
	NodeTypeStart       = domain.WorkflowNodeTypeStart
	NodeTypeEnd         = domain.WorkflowNodeTypeEnd
	NodeTypeCondition   = domain.WorkflowNodeTypeCondition
	NodeTypeBranchBlock = domain.WorkflowNodeTypeBranchBlock
	NodeTypeTryCatch    = domain.WorkflowNodeTypeTryCatch
	NodeTypeTryBlock    = domain.WorkflowNodeTypeTryBlock
	NodeTypeCatchBlock  = domain.WorkflowNodeTypeCatchBlock
	NodeTypeBizApply    = domain.WorkflowNodeTypeBizApply
	NodeTypeBizUpload   = domain.WorkflowNodeTypeBizUpload
	NodeTypeBizMonitor  = domain.WorkflowNodeTypeBizMonitor
	NodeTypeBizDeploy   = domain.WorkflowNodeTypeBizDeploy
	NodeTypeBizNotify   = domain.WorkflowNodeTypeBizNotify
)

type NodeIOEntry struct {
	Scope     string // 零值时表示全局的，否则表示指定节点的
	Type      string // 仅表示输入输出有值，表示变量无值
	Key       string
	Value     any
	ValueType string `options:"string | number | boolean"`
}

type Graph = domain.WorkflowGraph
