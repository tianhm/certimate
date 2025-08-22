package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/certimate-go/certimate/internal/domain/expr"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

const CollectionNameWorkflow = "workflow"

type Workflow struct {
	Meta
	Name          string                `json:"name" db:"name"`
	Description   string                `json:"description" db:"description"`
	Trigger       WorkflowTriggerType   `json:"trigger" db:"trigger"`
	TriggerCron   string                `json:"triggerCron" db:"triggerCron"`
	Enabled       bool                  `json:"enabled" db:"enabled"`
	GraphDraft    *WorkflowGraph        `json:"graphDraft" db:"graphDraft"`
	GraphContent  *WorkflowGraph        `json:"graphContent" db:"graphContent"`
	HasDraft      bool                  `json:"hasDraft" db:"hasDraft"`
	HasContent    bool                  `json:"hasContent" db:"hasContent"`
	LastRunId     string                `json:"lastRunId" db:"lastRunRef"`
	LastRunStatus WorkflowRunStatusType `json:"lastRunStatus" db:"lastRunStatus"`
	LastRunTime   time.Time             `json:"lastRunTime" db:"lastRunTime"`
}

type WorkflowGraph struct {
	Nodes []*WorkflowNode `json:"nodes"`
}

func (g *WorkflowGraph) GetNodeById(nodeId string) (*WorkflowNode, bool) {
	return g.getNodeInBlocksById(g.Nodes, nodeId)
}

func (g *WorkflowGraph) getNodeInBlocksById(blocks []*WorkflowNode, nodeId string) (*WorkflowNode, bool) {
	for _, node := range blocks {
		if node.Id == nodeId {
			return node, true
		}

		if len(node.Blocks) > 0 {
			if found, ok := g.getNodeInBlocksById(node.Blocks, nodeId); ok {
				return found, true
			}
		}
	}

	return nil, false
}

func (g *WorkflowGraph) Verify() error {
	if len(g.Nodes) < 2 {
		return fmt.Errorf("invalid nodes length of graph")
	} else if g.Nodes[0].Type != WorkflowNodeTypeStart {
		return fmt.Errorf("the first node is not a start node")
	} else if g.Nodes[len(g.Nodes)-1].Type != WorkflowNodeTypeEnd {
		return fmt.Errorf("the last node is not an end node")
	}

	return nil
}

func (g *WorkflowGraph) Clone() *WorkflowGraph {
	return &WorkflowGraph{
		Nodes: g.Nodes,
	}
}

type WorkflowTriggerType string

const (
	WorkflowTriggerTypeScheduled = WorkflowTriggerType("scheduled")
	WorkflowTriggerTypeManual    = WorkflowTriggerType("manual")
)

type WorkflowNode struct {
	Id     string           `json:"id"` // 节点 ID 只在该工作流中唯一，在全局中不保证唯一性
	Type   WorkflowNodeType `json:"type"`
	Data   WorkflowNodeData `json:"data"`
	Blocks []*WorkflowNode  `json:"blocks,omitempty"`
}

type WorkflowNodeType string

const (
	WorkflowNodeTypeStart       = WorkflowNodeType("start")
	WorkflowNodeTypeEnd         = WorkflowNodeType("end")
	WorkflowNodeTypeCondition   = WorkflowNodeType("condition")
	WorkflowNodeTypeBranchBlock = WorkflowNodeType("branchBlock")
	WorkflowNodeTypeTryCatch    = WorkflowNodeType("tryCatch")
	WorkflowNodeTypeTryBlock    = WorkflowNodeType("tryBlock")
	WorkflowNodeTypeCatchBlock  = WorkflowNodeType("catchBlock")
	WorkflowNodeTypeBizApply    = WorkflowNodeType("bizApply")
	WorkflowNodeTypeBizUpload   = WorkflowNodeType("bizUpload")
	WorkflowNodeTypeBizMonitor  = WorkflowNodeType("bizMonitor")
	WorkflowNodeTypeBizDeploy   = WorkflowNodeType("bizDeploy")
	WorkflowNodeTypeBizNotify   = WorkflowNodeType("bizNotify")
)

type WorkflowNodeData struct {
	Name   string             `json:"name"`
	Config WorkflowNodeConfig `json:"config"`
}

type WorkflowNodeConfig map[string]any

func (c WorkflowNodeConfig) AsBizApply() WorkflowNodeConfigForBizApply {
	return WorkflowNodeConfigForBizApply{
		Domains:               xmaps.GetString(c, "domains"),
		ContactEmail:          xmaps.GetString(c, "contactEmail"),
		ChallengeType:         xmaps.GetString(c, "challengeType"),
		Provider:              xmaps.GetString(c, "provider"),
		ProviderAccessId:      xmaps.GetString(c, "providerAccessId"),
		ProviderConfig:        xmaps.GetKVMapAny(c, "providerConfig"),
		KeyAlgorithm:          xmaps.GetOrDefaultString(c, "keyAlgorithm", string(CertificateKeyAlgorithmTypeRSA2048)),
		CAProvider:            xmaps.GetString(c, "caProvider"),
		CAProviderAccessId:    xmaps.GetString(c, "caProviderAccessId"),
		CAProviderConfig:      xmaps.GetKVMapAny(c, "caProviderConfig"),
		ACMEProfile:           xmaps.GetString(c, "acmeProfile"),
		Nameservers:           xmaps.GetString(c, "nameservers"),
		DnsPropagationWait:    xmaps.GetInt32(c, "dnsPropagationWait"),
		DnsPropagationTimeout: xmaps.GetInt32(c, "dnsPropagationTimeout"),
		DnsTTL:                xmaps.GetInt32(c, "dnsTTL"),
		DisableFollowCNAME:    xmaps.GetBool(c, "disableFollowCNAME"),
		DisableARI:            xmaps.GetBool(c, "disableARI"),
		SkipBeforeExpiryDays:  xmaps.GetInt32(c, "skipBeforeExpiryDays"),
	}
}

func (c WorkflowNodeConfig) AsBizUpload() WorkflowNodeConfigForBizUpload {
	return WorkflowNodeConfigForBizUpload{
		Certificate: xmaps.GetString(c, "certificate"),
		PrivateKey:  xmaps.GetString(c, "privateKey"),
		Domains:     xmaps.GetString(c, "domains"),
	}
}

func (c WorkflowNodeConfig) AsBizMonitor() WorkflowNodeConfigForBizMonitor {
	host := xmaps.GetString(c, "host")
	return WorkflowNodeConfigForBizMonitor{
		Host:        host,
		Port:        xmaps.GetOrDefaultInt32(c, "port", 443),
		Domain:      xmaps.GetOrDefaultString(c, "domain", host),
		RequestPath: xmaps.GetString(c, "path"),
	}
}

func (c WorkflowNodeConfig) AsBizDeploy() WorkflowNodeConfigForBizDeploy {
	return WorkflowNodeConfigForBizDeploy{
		CertificateOutputNodeId: xmaps.GetString(c, "certificateOutputNodeId"),
		Provider:                xmaps.GetString(c, "provider"),
		ProviderAccessId:        xmaps.GetString(c, "providerAccessId"),
		ProviderConfig:          xmaps.GetKVMapAny(c, "providerConfig"),
		SkipOnLastSucceeded:     xmaps.GetBool(c, "skipOnLastSucceeded"),
	}
}

func (c WorkflowNodeConfig) AsBizNotify() WorkflowNodeConfigForBizNotify {
	return WorkflowNodeConfigForBizNotify{
		Provider:             xmaps.GetString(c, "provider"),
		ProviderAccessId:     xmaps.GetString(c, "providerAccessId"),
		ProviderConfig:       xmaps.GetKVMapAny(c, "providerConfig"),
		Subject:              xmaps.GetString(c, "subject"),
		Message:              xmaps.GetString(c, "message"),
		SkipOnAllPrevSkipped: xmaps.GetBool(c, "skipOnAllPrevSkipped"),
	}
}

func (c WorkflowNodeConfig) AsBranchBlock() WorkflowNodeConfigForBranchBlock {
	expression := c["expression"]
	if expression == nil {
		return WorkflowNodeConfigForBranchBlock{}
	}

	exprRaw, _ := json.Marshal(expression)
	expr, err := expr.UnmarshalExpr([]byte(exprRaw))
	if err != nil {
		return WorkflowNodeConfigForBranchBlock{}
	}

	return WorkflowNodeConfigForBranchBlock{
		Expression: expr,
	}
}

type WorkflowNodeConfigForBizApply struct {
	Domains               string         `json:"domains"`                         // 域名列表，以半角分号分隔
	ContactEmail          string         `json:"contactEmail"`                    // 联系邮箱
	ChallengeType         string         `json:"challengeType"`                   // 验证方式。目前仅支持 dns-01
	Provider              string         `json:"provider"`                        // DNS 提供商
	ProviderAccessId      string         `json:"providerAccessId"`                // DNS 提供商授权记录 ID
	ProviderConfig        map[string]any `json:"providerConfig,omitempty"`        // DNS 提供商额外配置
	KeyAlgorithm          string         `json:"keyAlgorithm,omitempty"`          // 证书算法
	CAProvider            string         `json:"caProvider,omitempty"`            // CA 提供商（零值时使用全局配置）
	CAProviderAccessId    string         `json:"caProviderAccessId,omitempty"`    // CA 提供商授权记录 ID
	CAProviderConfig      map[string]any `json:"caProviderConfig,omitempty"`      // CA 提供商额外配置
	ACMEProfile           string         `json:"acmeProfile,omitempty"`           // ACME Profiles Extension
	Nameservers           string         `json:"nameservers,omitempty"`           // DNS 服务器列表，以半角分号分隔
	DnsPropagationWait    int32          `json:"dnsPropagationWait,omitempty"`    // DNS 传播等待时间，等同于 lego 的 `--dns-propagation-wait` 参数
	DnsPropagationTimeout int32          `json:"dnsPropagationTimeout,omitempty"` // DNS 传播检查超时时间（零值时使用提供商的默认值）
	DnsTTL                int32          `json:"dnsTTL,omitempty"`                // DNS 解析记录 TTL（零值时使用提供商的默认值）
	DisableFollowCNAME    bool           `json:"disableFollowCNAME,omitempty"`    // 是否关闭 CNAME 跟随
	DisableARI            bool           `json:"disableARI,omitempty"`            // 是否关闭 ARI
	SkipBeforeExpiryDays  int32          `json:"skipBeforeExpiryDays,omitempty"`  // 证书到期前多少天前跳过续期
}

type WorkflowNodeConfigForBizUpload struct {
	Certificate string `json:"certificate"` // 证书 PEM 内容
	PrivateKey  string `json:"privateKey"`  // 私钥 PEM 内容
	Domains     string `json:"domains,omitempty"`
}

type WorkflowNodeConfigForBizMonitor struct {
	Host        string `json:"host"`                  // 主机地址
	Port        int32  `json:"port,omitempty"`        // 端口（零值时默认值 443）
	Domain      string `json:"domain,omitempty"`      // 域名（零值时默认值 [Host]）
	RequestPath string `json:"requestPath,omitempty"` // 请求路径
}

type WorkflowNodeConfigForBizDeploy struct {
	CertificateOutputNodeId string         `json:"certificateOutputNodeId"`    // 前序证书输出节点 ID
	Provider                string         `json:"provider"`                   // 主机提供商
	ProviderAccessId        string         `json:"providerAccessId,omitempty"` // 主机提供商授权记录 ID
	ProviderConfig          map[string]any `json:"providerConfig,omitempty"`   // 主机提供商额外配置
	SkipOnLastSucceeded     bool           `json:"skipOnLastSucceeded"`        // 上次部署成功时是否跳过
}

type WorkflowNodeConfigForBizNotify struct {
	Provider             string         `json:"provider"`                 // 通知提供商
	ProviderAccessId     string         `json:"providerAccessId"`         // 通知提供商授权记录 ID
	ProviderConfig       map[string]any `json:"providerConfig,omitempty"` // 通知提供商额外配置
	Subject              string         `json:"subject"`                  // 通知主题
	Message              string         `json:"message"`                  // 通知内容
	SkipOnAllPrevSkipped bool           `json:"skipOnAllPrevSkipped"`     // 前序节点均已跳过时是否跳过
}

type WorkflowNodeConfigForBranchBlock struct {
	Expression expr.Expr `json:"expression"` // 条件表达式
}
