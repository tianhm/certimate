import { BizApplyNodeRegistry } from "./BizApplyNodeRegistry";
import { BizDeployNodeRegistry } from "./BizDeployNodeRegistry";
import { BizMonitorNodeRegistry } from "./BizMonitorNodeRegistry";
import { BizNotifyNodeRegistry } from "./BizNotifyNodeRegistry";
import { BizUploadNodeRegistry } from "./BizUploadNodeRegistry";
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

export { duplicateNodeJSON } from "./_shared";

export type * from "./typings";
