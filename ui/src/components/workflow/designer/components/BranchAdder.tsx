import { useTranslation } from "react-i18next";
import { FlowNodeBaseType, type FlowNodeEntity, type AdderProps as FlowgramAdderProps, useClientContext } from "@flowgram.ai/fixed-layout-editor";
import { Button } from "antd";
import { nanoid } from "nanoid";

import { BranchBlockNodeRegistry } from "../nodes/ConditionNode";
import { CatchBlockNodeRegistry } from "../nodes/TryCatchNode";
import { NodeType } from "../nodes/typings";

export interface BranchAdderProps extends FlowgramAdderProps {}

const BranchAdder = ({ node }: BranchAdderProps) => {
  const { t } = useTranslation();

  const ctx = useClientContext();
  const { operation, playground } = ctx;

  const handleAddBranch = () => {
    let block: FlowNodeEntity;
    switch (node.flowNodeType) {
      case NodeType.Condition:
        block = operation.addBlock(node, BranchBlockNodeRegistry.onAdd!(ctx, node));
        break;

      case NodeType.TryCatch:
        block = operation.addBlock(node, CatchBlockNodeRegistry.onAdd!(ctx, node));
        break;

      default:
        console.warn(`[certimate] unsupported node type for adding branch: '${node.flowNodeType}'`);
        block = operation.addBlock(node, {
          id: nanoid(),
          type: FlowNodeBaseType.BLOCK,
          data: {
            name: t("workflow_node.branch_block.default_name"),
          },
        });
        break;
    }

    setTimeout(() => {
      playground.scrollToView({
        bounds: block.bounds,
        scrollToCenter: true,
      });
    }, 1);
  };

  return playground.config.readonlyOrDisabled ? null : (
    <div
      className="relative"
      onMouseEnter={() => node.firstChild?.renderData?.toggleMouseEnter()}
      onMouseLeave={() => node.firstChild?.renderData?.toggleMouseLeave()}
    >
      <Button shape="round" size="small" onClick={handleAddBranch}>
        <span className="text-xs">{t("workflow.detail.design.nodes.add_branch")}</span>
      </Button>
    </div>
  );
};

export default BranchAdder;
