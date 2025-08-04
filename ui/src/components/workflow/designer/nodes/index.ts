import { BizApplyNodeRegistry, BizDeployNodeRegistry, BizMonitorNodeRegistry, BizNotifyNodeRegistry, BizUploadNodeRegistry } from "./BusinessNode";
import { BranchBlockNodeRegistry, ConditionNodeRegistry } from "./ConditionNode";
import { EndNodeRegistry } from "./EndNode";
import { StartNodeRegistry } from "./StartNode";
import { CatchBlockNodeRegistry, TryCatchNodeRegistry } from "./TryCatchNode";

export const getFlowNodeRegistries = () => {
  return [
    StartNodeRegistry,
    BizApplyNodeRegistry,
    BizUploadNodeRegistry,
    BizMonitorNodeRegistry,
    BizDeployNodeRegistry,
    BizNotifyNodeRegistry,
    ConditionNodeRegistry,
    BranchBlockNodeRegistry,
    TryCatchNodeRegistry,
    CatchBlockNodeRegistry,
    EndNodeRegistry,
  ];
};

export type * from "./typings";
