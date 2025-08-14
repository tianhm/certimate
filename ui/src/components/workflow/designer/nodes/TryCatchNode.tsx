import { getI18n } from "react-i18next";
import { Field } from "@flowgram.ai/fixed-layout-editor";
import { IconArrowsSplit, IconCircleX } from "@tabler/icons-react";

import { BaseNode, BranchNode } from "./_shared";
import { newNodeId } from "../_util";
import { NodeKindType, type NodeRegistry, NodeType } from "./typings";

export const TryCatchNodeRegistry: NodeRegistry = {
  type: NodeType.TryCatch,
  kindType: NodeKindType.Logic,

  meta: {
    helpText: getI18n().t("workflow_node.try_catch.help"),
    labelText: getI18n().t("workflow_node.try_catch.label"),

    icon: IconArrowsSplit,
    iconColor: "#fff",
    iconBgColor: "#373c43",

    clickable: false,
    expandable: false,
  },

  formMeta: {
    render: () => {
      return <BaseNode />;
    },
  },

  onAdd() {
    const { t } = getI18n();

    return {
      id: newNodeId(),
      type: NodeType.TryCatch,
      data: {
        name: t("workflow_node.try_catch.default_name"),
      },
      blocks: [
        {
          id: newNodeId(),
          type: NodeType.TryBlock,
          blocks: [],
        },
        {
          id: newNodeId(),
          type: NodeType.CatchBlock,
          blocks: [
            {
              id: newNodeId(),
              type: NodeType.End,
              data: {
                name: t("workflow_node.end.default_name"),
              },
            },
          ],
          data: {
            name: t("workflow_node.catch_block.default_name"),
          },
        },
      ],
    };
  },
};

export const CatchBlockNodeRegistry: NodeRegistry = {
  type: NodeType.CatchBlock,
  kindType: NodeKindType.Logic,

  meta: {
    labelText: getI18n().t("workflow_node.catch_block.label"),

    clickable: false,
    draggable: false,

    addDisable: true,
  },

  formMeta: {
    render: () => {
      return (
        <BranchNode
          description={
            <div className="flex items-center justify-center gap-2">
              <div className="flex items-center justify-center">
                <IconCircleX color="var(--color-error)" size="1.25em" stroke="1.25" />
              </div>
              <div className="truncate">
                <Field<string> name="name">{({ field: { value } }) => <>{value || "\u00A0"}</>}</Field>
              </div>
            </div>
          }
        />
      );
    },
  },

  canAdd: () => false,

  canDelete: (_, node) => {
    return node.parent != null && node.parent.blocks.length >= 2;
  },

  onAdd() {
    const { t } = getI18n();

    return {
      id: newNodeId(),
      type: NodeType.CatchBlock,
      blocks: [
        {
          id: newNodeId(),
          type: NodeType.End,
          data: {
            name: t("workflow_node.end.default_name"),
          },
        },
      ],
      data: {
        name: t("workflow_node.catch_block.default_name"),
      },
    };
  },
};
