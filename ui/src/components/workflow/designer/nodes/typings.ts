import {
  type FixedLayoutPluginContext,
  type FlowNodeEntity,
  type FlowNodeJSON,
  type FlowNodeMeta,
  type FlowNodeRegistry,
  type FormMeta,
  type FormRenderProps,
} from "@flowgram.ai/fixed-layout-editor";

import { WORKFLOW_NODE_TYPES, type WorkflowNode } from "@/domain/workflow";

export enum NodeType {
  Start = "start",
  End = "end",
  Condition = "condition",
  BranchBlock = "branchBlock",
  TryCatch = "tryCatch",
  TryBlock = "tryBlock",
  CatchBlock = "catchBlock",
  BizApply = "bizApply",
  BizUpload = "bizUpload",
  BizMonitor = "bizMonitor",
  BizDeploy = "bizDeploy",
  BizNotify = "bizNotify",
}

/* TYPE GUARD, PLEASE DO NOT REMOVE THESE! */
console.assert(NodeType.Start === WORKFLOW_NODE_TYPES.START);
console.assert(NodeType.End === WORKFLOW_NODE_TYPES.END);
console.assert(NodeType.Condition === WORKFLOW_NODE_TYPES.CONDITION);
console.assert(NodeType.BranchBlock === WORKFLOW_NODE_TYPES.BRANCHBLOCK);
console.assert(NodeType.TryCatch === WORKFLOW_NODE_TYPES.TRYCATCH);
console.assert(NodeType.TryBlock === WORKFLOW_NODE_TYPES.TRYBLOCK);
console.assert(NodeType.CatchBlock === WORKFLOW_NODE_TYPES.CATCHBLOCK);
console.assert(NodeType.BizApply === WORKFLOW_NODE_TYPES.BIZ_APPLY);
console.assert(NodeType.BizUpload === WORKFLOW_NODE_TYPES.BIZ_UPLOAD);
console.assert(NodeType.BizMonitor === WORKFLOW_NODE_TYPES.BIZ_MONITOR);
console.assert(NodeType.BizDeploy === WORKFLOW_NODE_TYPES.BIZ_DEPLOY);
console.assert(NodeType.BizNotify === WORKFLOW_NODE_TYPES.BIZ_NOTIFY);

export enum NodeKindType {
  Basis = "basis",
  Business = "business",
  Logic = "logic",
}

export interface NodeJSON extends FlowNodeJSON {
  data: WorkflowNode["data"] & {
    [key: string]: any;
  };
}

export interface DocumentJSON {
  nodes: NodeJSON[];
}

export interface NodeMeta extends FlowNodeMeta {
  style?: React.CSSProperties;
  helpText?: React.ReactNode;
  labelText?: React.ReactNode;
  icon?: React.ExoticComponent<any> | React.ComponentType<any>;
  iconColor?: string;
  iconBgColor?: string;
  clickable?: boolean;
}

export interface NodeRegistry<V extends NodeJSON["data"] = NodeJSON["data"]> extends FlowNodeRegistry<NodeMeta> {
  kind?: NodeKindType;
  formMeta?: Omit<FormMeta<V>, "render"> & {
    render: (props: FormRenderProps<V>) => React.ReactElement;
  };
  canAdd?: (ctx: FixedLayoutPluginContext, from: FlowNodeEntity) => boolean;
  canDelete?: (ctx: FixedLayoutPluginContext, from: FlowNodeEntity) => boolean;
  onAdd?: (ctx: FixedLayoutPluginContext, from: FlowNodeEntity) => FlowNodeJSON;
}
