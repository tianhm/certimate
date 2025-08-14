import { FlowNodeBaseType, type FlowNodeEntity } from "@flowgram.ai/fixed-layout-editor";
import { Immer } from "immer";
import { nanoid } from "nanoid";

import { type NodeJSON, NodeType } from "./nodes/typings";

/**
 * 返回一个新的节点 ID。
 * @returns {String}
 */
export const newNodeId = () => nanoid();

/**
 * 克隆节点 JSON 对象。节点及其子节点 ID 均会重新分配。
 * @param {NodeJSON} node
 * @param {Object} options
 * @returns {NodeJSON}
 */
export const duplicateNodeJSON = (node: NodeJSON, options?: { withCopySuffix?: boolean }) => {
  const { produce } = new Immer({ autoFreeze: false });
  const deepClone = (node: NodeJSON, { withCopySuffix, nodeIdMap }: { withCopySuffix: boolean; nodeIdMap: Map<string, string> }) => {
    return produce(node, (draft) => {
      draft.data ??= {};
      draft.id = newNodeId();
      draft.data.name = withCopySuffix ? `${draft.data?.name || ""}-copy` : `${draft.data?.name || ""}`;

      nodeIdMap.set(node.id, draft.id); // 原节点 ID 映射到新节点 ID

      if (draft.blocks) {
        draft.blocks = draft.blocks.map((block) => deepClone(block as NodeJSON, { withCopySuffix: false, nodeIdMap }));
      }

      if (draft.data?.config) {
        switch (draft.type) {
          case NodeType.BizDeploy:
            {
              const prevNodeId = draft.data.config.certificate?.split("#")?.[0];
              if (nodeIdMap.has(prevNodeId)) {
                draft.data.config = {
                  ...draft.data.config,
                  certificate: `${nodeIdMap.get(prevNodeId)}#certificate`,
                };
              }
            }
            break;

          case NodeType.Condition:
            {
              const stack = [] as any[];
              const expr = draft.data.config.expression;
              if (expr) {
                stack.push(expr);
                while (stack.length > 0) {
                  const n = stack.pop()!;
                  if ("left" in n) {
                    stack.push(n.left);
                    if ("selector" in n.left) {
                      const prevNodeId = n.left.selector.id;
                      if (nodeIdMap.has(prevNodeId)) {
                        n.left.selector.id = nodeIdMap.get(prevNodeId)!;
                      }
                    }
                  }
                  if ("right" in n) {
                    stack.push(n.right);
                  }
                }
                draft.data.config = {
                  ...draft.data.config,
                  expression: expr,
                };
              }
            }
            break;
        }
      }

      return draft;
    });
  };

  return deepClone(node, { withCopySuffix: options?.withCopySuffix ?? true, nodeIdMap: new Map() });
};

/**
 * 获取指定节点到根节点为止的所有前序节点。不包括自身和根节点或开始节点。
 * @param {FlowNodeEntity} node
 * @returns {FlowNodeEntity[]}
 */
export const getAllPreviousNodes = (node: FlowNodeEntity): FlowNodeEntity[] => {
  if (node == null) return [];

  // TODO: 不应该获取到旁路分支
  // // 先获取单一链路（即不包含分支）的全部节点
  // const chains: FlowNodeEntity[] = [];
  // let chain: FlowNodeEntity | undefined = node;
  // while (chain) {
  //   if (chain.isStart || chain.flowNodeType === FlowNodeBaseType.ROOT) {
  //     break;
  //   }

  //   chains.push(chain);
  //   chain = chain.pre ?? chain.parent;
  // }

  // 再获取实际的全部节点
  const visited = new Set<string>();
  const result: FlowNodeEntity[] = [];
  let current: FlowNodeEntity | undefined = node;
  while (current) {
    if (current.isStart || current.flowNodeType === FlowNodeBaseType.ROOT) {
      break;
    }

    if (current.flowNodeType === NodeType.Condition) {
      /**
       * condition
       *   blockIcon
       *   inlineBlocks
       *     branchBlock_1
       *       blockOrderIcon
       *       ...
       *     branchBlock_2
       *       blockOrderIcon
       *       ...
       */
      current.lastBlock?.blocks?.forEach((block) => {
        block.allChildren?.forEach((child) => {
          if (!visited.has(child.id)) {
            visited.add(child.id);
            result.push(child);
          }
        });
      });
    } else if (current.flowNodeType === NodeType.TryCatch) {
      /**
       * tryCatch
       *   blockIcon
       *   mainInlineBlocks
       *     tryBlock
       *       trySlot
       *       ...
       *     catchInlineBlocks
       *       catchBlock_1
       *         blockOrderIcon
       *         ...
       *         end
       *       catchBlock_2
       *         blockOrderIcon
       *         ...
       *         end
       */
      current.lastBlock?.blocks?.forEach((block) => {
        block.allChildren?.forEach((child) => {
          if (!visited.has(child.id)) {
            visited.add(child.id);
            result.push(child);
          }
        });
      });
    }

    if (!visited.has(current.id)) {
      visited.add(current.id);
      result.push(current);
    }
    current = current.pre ?? current.parent;
  }

  const prevNodes = result.filter((e) => {
    if (e.id === node.id) return false;
    if (e.isTypeOrExtendType(FlowNodeBaseType.BLOCK_ICON)) return false;
    if (e.isTypeOrExtendType(FlowNodeBaseType.BLOCK_ORDER_ICON)) return false;
    if (e.isTypeOrExtendType(FlowNodeBaseType.INLINE_BLOCKS)) return false;
    if (e.isTypeOrExtendType("trySlot")) return false;

    return true;
  });
  console.log(node.document.root);
  console.log(prevNodes);
  return prevNodes;
};
