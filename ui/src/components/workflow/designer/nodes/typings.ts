import { FlowNodeBaseType } from "@flowgram.ai/document";
import {
  type FixedLayoutPluginContext,
  type FlowNodeEntity,
  type FlowNodeJSON,
  type FlowNodeMeta,
  type FlowNodeRegistry,
  type FormMeta,
  type FormRenderProps,
} from "@flowgram.ai/fixed-layout-editor";

export enum NodeType {
  Start = FlowNodeBaseType.START,
  End = FlowNodeBaseType.END,
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

export interface NodeJSON extends FlowNodeJSON {
  data: {
    name?: string;
    disabled?: boolean;
    [key: string]: any;
  };
}

export interface DocumentJSON {
  nodes: NodeJSON[];
}

export interface NodeMeta extends FlowNodeMeta {
  style?: React.CSSProperties;
  helpText?: React.ReactNode;
  icon?: React.Component;
  iconColor?: string;
  iconBgColor?: string;
}

export interface NodeRegistry<V extends NodeJSON["data"] = NodeJSON["data"]> extends FlowNodeRegistry<NodeMeta> {
  meta?: FlowNodeMeta;
  formMeta?: Omit<FormMeta<V>, "render"> & {
    render: (props: FormRenderProps<V>) => React.ReactElement;
  };
  canAdd?: (ctx: FixedLayoutPluginContext, from: FlowNodeEntity) => boolean;
  canDelete?: (ctx: FixedLayoutPluginContext, from: FlowNodeEntity) => boolean;
  onAdd?: (ctx: FixedLayoutPluginContext, from: FlowNodeEntity) => FlowNodeJSON;
}
