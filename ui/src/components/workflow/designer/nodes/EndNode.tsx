import { getI18n } from "react-i18next";
import { FlowNodeBaseType } from "@flowgram.ai/fixed-layout-editor";
import { IconLogout } from "@tabler/icons-react";
import { nanoid } from "nanoid";

import { BaseNode } from "./_shared";
import { type NodeRegistry, NodeType } from "./typings";

export const EndNodeRegistry: NodeRegistry = {
  type: NodeType.End,

  meta: {
    helpText: getI18n().t("workflow_node.end.help"),
    labelText: getI18n().t("workflow_node.end.label"),

    icon: IconLogout,
    iconColor: "#fff",
    iconBgColor: "#336df4",

    isNodeEnd: true,

    clickable: false,
    expandable: false,
    selectable: false,

    copyDisable: true,
  },

  formMeta: {
    render: () => {
      return <BaseNode></BaseNode>;
    },
  },

  canAdd(_, from) {
    // You can only add to the last node of the branch
    if (!from.isLast) return false;

    /**
     * condition
     *  blockIcon
     *  inlineBlocks
     *    block1
     *      blockOrderIcon
     *      <---- [add end]
     *    block2
     *      blockOrderIcon
     *      end
     */
    // originParent can determine whether it is condition , and then determine whether it is the last one
    // https://github.com/bytedance/flowgram.ai/pull/146
    if (from.parent && from.parent.parent?.flowNodeType === FlowNodeBaseType.INLINE_BLOCKS && from.parent.originParent && !from.parent.originParent.isLast) {
      const allBranches = from.parent.parent!.blocks;
      // Determine whether the last node of all branch is end, All branches are not allowed to be end
      const branchEndCount = allBranches.filter((block) => block.blocks[block.blocks.length - 1]?.getNodeMeta().isNodeEnd).length;
      return branchEndCount < allBranches.length - 1;
    }

    return true;
  },

  canDelete(ctx, node) {
    return node.parent !== ctx.document.root;
  },

  onAdd() {
    const { t } = getI18n();

    return {
      id: nanoid(),
      type: NodeType.End,
      data: {
        name: t("workflow_node.end.default_name"),
      },
    };
  },
};
