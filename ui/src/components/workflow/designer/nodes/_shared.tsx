import { useContext, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { type FlowNodeEntity, useClientContext } from "@flowgram.ai/fixed-layout-editor";
import { IconCopy, IconDotsVertical, IconGripVertical, IconLabel, IconX } from "@tabler/icons-react";
import { Button, type ButtonProps, Card, Dropdown, Popover, Tooltip, Typography, theme } from "antd";
import { Immer } from "immer";
import { nanoid } from "nanoid";

import { mergeCls } from "@/utils/css";

import { type NodeJSON, type NodeRegistry, NodeType } from "./typings";
import { NodeRenderContext } from "../NodeRenderContext";

export const BaseNode = ({ className, style, children }: { className?: string; style?: React.CSSProperties; children?: React.ReactNode }) => {
  const { token: themeToken } = theme.useToken();

  const ctx = useClientContext();
  const { playground } = ctx;

  const nodeRender = useContext(NodeRenderContext);
  const nodeRegistry = nodeRender.node.getNodeRegistry<NodeRegistry>();

  const NodeIcon = nodeRegistry.meta?.icon;
  const renderNodeIcon = () => {
    return NodeIcon == null ? null : (
      <div
        className="mr-2 flex size-9 items-center justify-center rounded-lg bg-white text-primary shadow-md dark:bg-stone-200"
        style={{
          color: nodeRegistry.meta?.iconColor,
          backgroundColor: nodeRegistry.meta?.iconBgColor,
        }}
      >
        <NodeIcon size="1.75em" color={nodeRegistry.meta?.iconColor} stroke="1.25" />
      </div>
    );
  };

  return (
    <div className={mergeCls("relative w-[320px] group/node", className)} style={style}>
      <Card className="rounded-xl shadow-sm" styles={{ body: { padding: 0 } }} hoverable>
        <div className="flex items-center gap-1 overflow-hidden p-3">
          {nodeRegistry.meta?.helpText == null ? (
            renderNodeIcon()
          ) : (
            <Tooltip title={<span dangerouslySetInnerHTML={{ __html: nodeRegistry.meta.helpText }}></span>} mouseEnterDelay={1}>
              <div className="cursor-help">{renderNodeIcon()}</div>
            </Tooltip>
          )}
          <div className="flex-1 overflow-hidden">
            <div className="truncate">
              <Typography.Text>{nodeRender.data?.name || "\u00A0"}</Typography.Text>
            </div>
            {children != null && (
              <div className="truncate text-xs" style={{ color: themeToken.colorTextTertiary }}>
                {children}
              </div>
            )}
          </div>
          <div className="-mr-2 ml-1" onClick={(e) => e.stopPropagation()}>
            <NodeMenuButton className="opacity-0 transition-opacity group-hover/node:opacity-100" size="small" />
          </div>
        </div>
      </Card>
      {!playground.config.readonlyOrDisabled && nodeRegistry.meta?.draggable === true && (
        <div className="absolute top-1/2 -left-4 hidden -translate-y-1/2 group-hover/node:block">
          <IconGripVertical size="1em" stroke="1" />
        </div>
      )}
    </div>
  );
};

export const BlockNode = ({ className, style, children }: { className?: string; style?: React.CSSProperties; children?: React.ReactNode }) => {
  const ctx = useClientContext();
  const { playground } = ctx;

  const nodeRender = useContext(NodeRenderContext);
  const nodeRegistry = nodeRender.node.getNodeRegistry<NodeRegistry>();

  return (
    <Popover classNames={{ root: "shadow-md" }} styles={{ body: { padding: 0 } }} arrow={false} content={<NodeMenuButton variant="text" />} placement="right">
      <div className={mergeCls("relative w-[240px] group/node", className)} style={style}>
        <Card className="rounded-xl shadow-sm" styles={{ body: { padding: 0 } }} hoverable>
          <div className="overflow-hidden px-3 py-2">
            <div className="truncate text-center">{children ?? (nodeRender.data?.name || "\u00A0")}</div>
          </div>
        </Card>
        {!playground.config.readonlyOrDisabled && nodeRegistry.meta?.draggable === true && (
          <div className="absolute top-1/2 -left-4 hidden -translate-y-1/2 group-hover/node:block">
            <IconGripVertical size="1em" stroke="1" />
          </div>
        )}
      </div>
    </Popover>
  );
};

export const NodeMenuButton = ({ className, style, ...props }: ButtonProps) => {
  const { t } = useTranslation();

  const ctx = useClientContext();
  const { operation, playground } = ctx;

  const { node, deleteNode, isBlockOrderIcon } = useContext(NodeRenderContext);
  const nodeRegistry = node.getNodeRegistry<NodeRegistry>();

  const getLatestDuplicateDisabledState = () => {
    if (nodeRegistry.meta?.copyDisable != null && nodeRegistry.meta.copyDisable) {
      return true;
    }
    return false;
  };
  const getLatestRemoveDisabledState = () => {
    if (nodeRegistry.meta?.deleteDisable != null && nodeRegistry.meta.deleteDisable) {
      return true;
    }
    if (nodeRegistry.canDelete != null) {
      return !nodeRegistry.canDelete(ctx, node);
    }
    return false;
  };
  const [menuDuplicateDisabled, setMenuDuplicateDisabled] = useState(() => getLatestDuplicateDisabledState());
  const [menuRemoveDisabled, setMenuRemoveDisabled] = useState(() => getLatestRemoveDisabledState());
  useEffect(() => {
    // 这里不能使用 useMemo() 来决定 menuRemoveDisabled，因为依赖项没有发生改变（对象引用始终是同一个）
    // 因此需要使用 useEffect() 来监听 node 和 node.parent 的变化，并更新 menuRemoveDisabled 的状态
    const callback = () => {
      setMenuDuplicateDisabled(getLatestDuplicateDisabledState());
      setMenuRemoveDisabled(getLatestRemoveDisabledState());
    };
    const disposable1 = node.onEntityChange(callback);
    const disposable2 = node.parent?.onEntityChange?.(callback);
    return () => {
      disposable1?.dispose();
      disposable2?.dispose();
    };
  }, []);

  const handleClickRename = () => {
    alert("TODO: rename");
  };

  const handleClickDuplicate = () => {
    if (menuDuplicateDisabled) {
      return;
    }

    const parent = node.originParent ?? node.parent;
    if (parent != null) {
      const nodeJSON = duplicateNodeJSON(node.toJSON() as NodeJSON);

      let block: FlowNodeEntity;
      if (isBlockOrderIcon) {
        block = operation.addBlock(parent, nodeJSON);
      } else {
        block = operation.addFromNode(node, nodeJSON);
      }

      setTimeout(() => {
        playground.scrollToView({
          bounds: block.bounds,
          scrollToCenter: true,
        });
      }, 1);
    }
  };

  const handleClickRemove = () => {
    if (menuRemoveDisabled) {
      return;
    }

    deleteNode();
  };

  return playground.config.readonlyOrDisabled ? null : (
    <Dropdown
      arrow={false}
      destroyOnHidden
      menu={{
        items: [
          {
            key: "rename",
            label: isBlockOrderIcon ? t("workflow.detail.design.nodes.rename_branch") : t("workflow.detail.design.nodes.rename_node"),
            icon: <IconLabel size="1em" />,
            onClick: handleClickRename,
          },
          {
            key: "duplicate",
            label: isBlockOrderIcon ? t("workflow.detail.design.nodes.duplicate_branch") : t("workflow.detail.design.nodes.duplicate_node"),
            icon: <IconCopy size="1em" />,
            disabled: menuDuplicateDisabled,
            onClick: handleClickDuplicate,
          },
          {
            type: "divider",
          },
          {
            key: "remove",
            label: isBlockOrderIcon ? t("workflow.detail.design.nodes.remove_branch") : t("workflow.detail.design.nodes.remove_node"),
            icon: <IconX size="1em" />,
            danger: true,
            disabled: menuRemoveDisabled,
            onClick: handleClickRemove,
          },
        ],
      }}
      trigger={["click"]}
    >
      <Button className={className} style={style} icon={<IconDotsVertical color="grey" size="1.25em" />} type="text" {...props} />
    </Dropdown>
  );
};

// TODO: 应放至领域层
export const duplicateNodeJSON = (node: NodeJSON, options?: { withCopySuffix?: boolean }) => {
  const { produce } = new Immer({ autoFreeze: false });
  const deepClone = (node: NodeJSON, { withCopySuffix, nodeIdMap }: { withCopySuffix: boolean; nodeIdMap: Map<string, string> }) => {
    return produce(node, (draft) => {
      draft.data ??= {};
      draft.id = nanoid();
      draft.data.name = withCopySuffix ? `${draft.data?.name || ""}-copy` : `${draft.data?.name || ""}`;
      delete draft.data.disabled;

      nodeIdMap.set(node.id, draft.id); // 原节点 ID 映射到新节点 ID

      if (draft.blocks) {
        draft.blocks = draft.blocks.map((block) => deepClone(block as NodeJSON, { withCopySuffix, nodeIdMap }));
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
