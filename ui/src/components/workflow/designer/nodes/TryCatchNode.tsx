import { getI18n } from "react-i18next";
import { ValidateTrigger } from "@flowgram.ai/fixed-layout-editor";
import { IconArrowsSplit, IconCircleX } from "@tabler/icons-react";
import { nanoid } from "nanoid";

import { BaseNode, BlockNode } from "./_shared";
import { type NodeRegistry, NodeType } from "./typings";

export const TryCatchNodeRegistry: NodeRegistry = {
  type: NodeType.TryCatch,

  meta: {
    helpText: getI18n().t("workflow_node.try_catch.help"),

    icon: IconArrowsSplit,
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
      type: NodeType.TryCatch,
      data: {
        name: t("workflow_node.try_catch.default_name"),
      },
      blocks: [
        {
          id: nanoid(),
          type: NodeType.TryBlock,
          blocks: [],
        },
        {
          id: nanoid(),
          type: NodeType.CatchBlock,
          blocks: [
            {
              id: nanoid(),
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

  meta: {
    addDisable: true,
    copyDisable: true,
  },

  formMeta: {
    validateTrigger: ValidateTrigger.onChange,
    render: ({ form }) => {
      return (
        <BlockNode>
          <div className="flex items-center justify-center gap-2">
            <div className="flex items-center justify-center">
              <IconCircleX color="var(--color-error)" size="1.25em" stroke="1.25" />
            </div>
            <div className="truncate">{form.getValueIn<string>("name") || "\u00A0"}</div>
          </div>
        </BlockNode>
      );
    },
  },

  canAdd: () => false,

  canDelete: (_, node) => node.parent != null && node.parent.blocks.length >= 2,

  onAdd() {
    const { t } = getI18n();

    return {
      id: nanoid(),
      type: NodeType.CatchBlock,
      blocks: [],
      data: {
        name: t("workflow_node.catch_block.default_name"),
      },
    };
  },
};
