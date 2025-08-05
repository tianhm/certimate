import { getI18n } from "react-i18next";
import { FlowNodeBaseType, FlowNodeSplitType, ValidateTrigger } from "@flowgram.ai/fixed-layout-editor";
import { IconFilter, IconFilterFilled, IconSitemap } from "@tabler/icons-react";
import { nanoid } from "nanoid";

import { BaseNode, BlockNode } from "./_shared";
import { type NodeRegistry, NodeType } from "./typings";

export const ConditionNodeRegistry: NodeRegistry = {
  type: NodeType.Condition,

  extend: FlowNodeSplitType.DYNAMIC_SPLIT,

  meta: {
    helpText: getI18n().t("workflow_node.condition.help"),

    icon: IconSitemap,
    iconColor: "#fff",
    iconBgColor: "#373c43",

    expandable: false,
  },

  formMeta: {
    validateTrigger: ValidateTrigger.onChange,
    render: () => {
      return <BaseNode></BaseNode>;
    },
  },

  onAdd() {
    const { t } = getI18n();

    return {
      id: nanoid(),
      type: NodeType.Condition,
      data: {
        name: t("workflow_node.condition.default_name"),
      },
      blocks: [
        {
          id: nanoid(),
          type: NodeType.BranchBlock,
          blocks: [],
          data: {
            name: t("workflow_node.branch_block.default_name") + " 1",
          },
        },
        {
          id: nanoid(),
          type: NodeType.BranchBlock,
          data: {
            name: t("workflow_node.branch_block.default_name") + " 2",
          },
        },
      ],
    };
  },
};

export const BranchBlockNodeRegistry: NodeRegistry = {
  type: NodeType.BranchBlock,

  extend: FlowNodeBaseType.BLOCK,

  meta: {
    addDisable: true,
    copyDisable: true,
  },

  formMeta: {
    validateTrigger: ValidateTrigger.onChange,
    render: ({ form }) => {
      const fieldExpr = form.getValueIn("config.expression");

      return (
        <BlockNode>
          <div className="flex items-center justify-center gap-2">
            <div className="flex items-center justify-center">
              {fieldExpr ? <IconFilterFilled color="var(--color-primary)" size="1.25em" stroke="1.25" /> : <IconFilter size="1.25em" stroke="1.25" />}
            </div>
            <div className="truncate">{form.getValueIn<string>("name") || "\u00A0"}</div>
          </div>
        </BlockNode>
      );
    },
  },

  canAdd: () => {
    return false;
  },

  canDelete: (_, node) => {
    return node.parent != null && node.parent.blocks.length >= 2;
  },

  onAdd() {
    const { t } = getI18n();

    return {
      id: nanoid(),
      type: NodeType.BranchBlock,
      data: {
        name: t("workflow_node.branch_block.default_name"),
      },
    };
  },
};
