import { getI18n } from "react-i18next";
import { Field, FlowNodeBaseType, FlowNodeSplitType } from "@flowgram.ai/fixed-layout-editor";
import { IconFilter, IconFilterFilled, IconSitemap } from "@tabler/icons-react";
import { Typography } from "antd";
import { nanoid } from "nanoid";

import { type Expr, ExprType } from "@/domain/workflow";

import { BaseNode, BranchNode } from "./_shared";
import { type NodeRegistry, NodeType } from "./typings";

export const ConditionNodeRegistry: NodeRegistry = {
  type: NodeType.Condition,

  extend: FlowNodeSplitType.DYNAMIC_SPLIT,

  meta: {
    helpText: getI18n().t("workflow_node.condition.help"),
    labelText: getI18n().t("workflow_node.condition.label"),

    icon: IconSitemap,
    iconColor: "#fff",
    iconBgColor: "#373c43",

    clickable: false,
    expandable: false,

    deleteDisable: false,
  },

  formMeta: {
    render: () => {
      return <BaseNode />;
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
    labelText: getI18n().t("workflow_node.branch_block.label"),

    icon: IconSitemap,
    iconColor: "#fff",
    iconBgColor: "#373c43",

    clickable: true,

    addDisable: true,
    copyDisable: true,
  },

  formMeta: {
    render: () => {
      const { t } = getI18n();

      return (
        <BranchNode
          description={
            <>
              <div className="flex items-center justify-center gap-2">
                <div className="flex items-center justify-center">
                  <Field<Expr> name="config.expression">
                    {({ field: { value } }) => (
                      <>
                        {value == null ? (
                          <IconFilter size="1.25em" stroke="1.25" />
                        ) : (
                          <IconFilterFilled color="var(--color-primary)" size="1.25em" stroke="1.25" />
                        )}
                      </>
                    )}
                  </Field>
                </div>
                <div className="truncate">
                  <Field<string> name="name">{({ field: { value } }) => <>{value || "\u00A0"}</>}</Field>
                </div>
              </div>
              <div className="mt-1">
                <div className="truncate">
                  <Field<Expr> name="config.expression">
                    {({ field: { value } }) => (
                      <Typography.Text className="text-xs" type="secondary">
                        {value == null
                          ? t("workflow_node.branch_block.state.no")
                          : value.type === ExprType.Logical && value.operator === "and"
                            ? t("workflow_node.branch_block.state.and")
                            : t("workflow_node.branch_block.state.or")}
                      </Typography.Text>
                    )}
                  </Field>
                </div>
              </div>
            </>
          }
        />
      );
    },
  },

  canAdd: () => {
    return false;
  },

  canDelete: (_, node) => {
    return node.parent != null && node.parent.blocks.length >= 2;
  },

  onAdd(_, from) {
    const { t } = getI18n();

    let nodeName = t("workflow_node.branch_block.default_name");
    if (from != null) {
      const siblingLength = from.blocks?.find((b) => b.isInlineBlocks)?.blocks?.length;
      if (siblingLength != null) {
        nodeName = `${nodeName} ${siblingLength + 1}`;
      }
    }

    return {
      id: nanoid(),
      type: NodeType.BranchBlock,
      data: {
        name: nodeName,
      },
    };
  },
};
