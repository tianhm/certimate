import { type ContextType, useContext, useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { Field, type FieldRenderProps, type FlowNodeEntity, useClientContext } from "@flowgram.ai/fixed-layout-editor";
import { IconCopy, IconDotsVertical, IconGripVertical, IconLabel, IconX } from "@tabler/icons-react";
import { Button, type ButtonProps, Card, Dropdown, Input, type InputRef, Popover, Tooltip, theme } from "antd";
import { Immer } from "immer";
import { nanoid } from "nanoid";

import { mergeCls } from "@/utils/css";

import { type NodeJSON, type NodeRegistry, NodeType } from "./typings";
import { NodeRenderContext } from "../NodeRenderContext";

const useNodeRenamingInput = ({ nodeRender }: { nodeRender: ContextType<typeof NodeRenderContext> }) => {
  const inputRef = useRef<InputRef>(null);
  const [inputVisible, setInputVisible] = useState(false);
  const [inputValue, setInputValue] = useState("");

  const handleNodeRenameClick = () => {
    setInputVisible(true);
    setInputValue(nodeRender.data?.name);
    setTimeout(() => {
      inputRef.current?.focus({ cursor: "end" });
    }, 0);
  };

  const handleNodeNameChange = (value: string) => {
    setInputValue(value);
  };

  const handleNodeNameConfirm = async (value: string) => {
    value = value.trim();
    if (!value || value === (nodeRender.data?.name || "")) {
      setInputVisible(false);
      return;
    }

    setInputVisible(false);

    nodeRender.updateData({ name: value });
  };

  return {
    inputRef: inputRef,
    visible: inputVisible,
    value: inputValue,
    onClick: handleNodeRenameClick,
    onChange: handleNodeNameChange,
    onConfirm: handleNodeNameConfirm,
  };
};

export interface BaseNodeProps {
  className?: string;
  style?: React.CSSProperties;
  children?: React.ReactNode;
}

export const BaseNode = ({ className, style, children }: BaseNodeProps) => {
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

  const {
    inputRef,
    visible: inputVisible,
    value: inputValue,
    onClick: handleNodeRenameClick,
    onChange: handleNodeNameChange,
    onConfirm: handleNodeNameConfirm,
  } = useNodeRenamingInput({ nodeRender });

  return (
    <div className={mergeCls("relative w-[320px] group/node", className)} style={style}>
      <Card className="rounded-xl shadow-sm" styles={{ body: { padding: 0 } }} hoverable>
        <div className={mergeCls("flex items-center gap-1 overflow-hidden p-3", inputVisible ? "invisible" : "visible")}>
          {nodeRegistry.meta?.helpText == null ? (
            renderNodeIcon()
          ) : (
            <Tooltip title={<span dangerouslySetInnerHTML={{ __html: nodeRegistry.meta.helpText }}></span>} mouseEnterDelay={1}>
              <div className="cursor-help">{renderNodeIcon()}</div>
            </Tooltip>
          )}
          <div className="flex-1 overflow-hidden">
            <div className="truncate">
              <Field name="name">{({ field: { value } }: FieldRenderProps<string>) => <>{value || "\u00A0"}</>}</Field>
            </div>
            {children != null && (
              <div className="truncate text-xs" style={{ color: themeToken.colorTextTertiary }}>
                {children}
              </div>
            )}
          </div>
          <div className="-mr-2 ml-1" onClick={(e) => e.stopPropagation()}>
            <NodeMenuButton
              className="opacity-0 transition-opacity group-hover/node:opacity-100"
              size="small"
              onMenuSelect={(key) => {
                switch (key) {
                  case "rename":
                    handleNodeRenameClick();
                    break;
                }
              }}
            />
          </div>
        </div>
      </Card>
      {!playground.config.readonlyOrDisabled && nodeRegistry.meta?.draggable === true && (
        <div className="absolute top-1/2 -left-4 hidden -translate-y-1/2 group-hover/node:block">
          <IconGripVertical size="1em" stroke="1" />
        </div>
      )}
      {!playground.config.readonlyOrDisabled && (
        <div className={mergeCls("absolute top-1/2 left-2 right-2 -translate-y-1/2", inputVisible ? "block" : "hidden")}>
          <Input
            ref={inputRef}
            maxLength={100}
            variant="underlined"
            value={inputValue}
            onBlur={(e) => handleNodeNameConfirm(e.target.value)}
            onChange={(e) => handleNodeNameChange(e.target.value)}
            onPressEnter={(e) => e.currentTarget.blur()}
          />
        </div>
      )}
    </div>
  );
};

export interface BranchLikeNodeProps extends BaseNodeProps {}

export const BranchLikeNode = ({ className, style, children }: BranchLikeNodeProps) => {
  const ctx = useClientContext();
  const { playground } = ctx;

  const nodeRender = useContext(NodeRenderContext);
  const nodeRegistry = nodeRender.node.getNodeRegistry<NodeRegistry>();

  const {
    inputRef,
    visible: inputVisible,
    value: inputValue,
    onClick: handleNodeRenameClick,
    onChange: handleNodeNameChange,
    onConfirm: handleNodeNameConfirm,
  } = useNodeRenamingInput({ nodeRender });

  return (
    <Popover
      classNames={{ root: "shadow-md" }}
      styles={{ body: { padding: 0 } }}
      arrow={false}
      content={
        inputVisible ? null : (
          <NodeMenuButton
            variant="text"
            onMenuSelect={(key) => {
              switch (key) {
                case "rename":
                  handleNodeRenameClick();
                  break;
              }
            }}
          />
        )
      }
      placement="right"
    >
      <div className={mergeCls("relative w-[240px] group/node", className)} style={style}>
        <Card className="rounded-xl shadow-sm" styles={{ body: { padding: 0 } }} hoverable>
          <div className={mergeCls("overflow-hidden px-3 py-2", inputVisible ? "invisible" : "visible")}>
            <div className="truncate text-center">
              {children ?? <Field name="name">{({ field: { value } }: FieldRenderProps<string>) => <>{value || "\u00A0"}</>}</Field>}
            </div>
          </div>
        </Card>
        {!playground.config.readonlyOrDisabled && nodeRegistry.meta?.draggable === true && (
          <div className="absolute top-1/2 -left-4 hidden -translate-y-1/2 group-hover/node:block">
            <IconGripVertical size="1em" stroke="1" />
          </div>
        )}
        {!playground.config.readonlyOrDisabled && (
          <div className={mergeCls("absolute top-1/2 left-2 right-2 -translate-y-1/2", inputVisible ? "block" : "hidden")}>
            <Input
              ref={inputRef}
              maxLength={100}
              variant="underlined"
              value={inputValue}
              onBlur={(e) => handleNodeNameConfirm(e.target.value)}
              onChange={(e) => handleNodeNameChange(e.target.value)}
              onPressEnter={(e) => e.currentTarget.blur()}
            />
          </div>
        )}
      </div>
    </Popover>
  );
};

export interface NodeMenuButtionProps extends ButtonProps {
  className?: string;
  style?: React.CSSProperties;
  onMenuSelect?: (key: "rename" | "duplicate" | "remove") => void;
}

export const NodeMenuButton = ({ className, style, onMenuSelect, ...props }: NodeMenuButtionProps) => {
  const { t } = useTranslation();

  const ctx = useClientContext();
  const { operation, playground } = ctx;

  const { node, ...nodeRender } = useContext(NodeRenderContext);
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
    onMenuSelect?.("rename");
  };

  const handleClickDuplicate = () => {
    if (menuDuplicateDisabled) {
      return;
    }

    const parent = node.originParent ?? node.parent;
    if (parent != null) {
      const nodeJSON = duplicateNodeJSON(node.toJSON() as NodeJSON);

      let block: FlowNodeEntity;
      if (nodeRender.isBlockOrderIcon) {
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

      onMenuSelect?.("duplicate");
    }
  };

  const handleClickRemove = () => {
    if (menuRemoveDisabled) {
      return;
    }

    nodeRender.deleteNode();

    onMenuSelect?.("remove");
  };

  return playground.config.readonlyOrDisabled ? null : (
    <Dropdown
      arrow={false}
      destroyOnHidden
      menu={{
        items: [
          {
            key: "rename",
            label: nodeRender.isBlockOrderIcon ? t("workflow.detail.design.nodes.rename_branch") : t("workflow.detail.design.nodes.rename_node"),
            icon: <IconLabel size="1em" />,
            onClick: handleClickRename,
          },
          {
            key: "duplicate",
            label: nodeRender.isBlockOrderIcon ? t("workflow.detail.design.nodes.duplicate_branch") : t("workflow.detail.design.nodes.duplicate_node"),
            icon: <IconCopy size="1em" />,
            disabled: menuDuplicateDisabled,
            onClick: handleClickDuplicate,
          },
          {
            type: "divider",
          },
          {
            key: "remove",
            label: nodeRender.isBlockOrderIcon ? t("workflow.detail.design.nodes.remove_branch") : t("workflow.detail.design.nodes.remove_node"),
            icon: <IconX size="1em" />,
            danger: true,
            disabled: menuRemoveDisabled,
            onClick: handleClickRemove,
          },
        ],
      }}
      overlayStyle={{
        zIndex: 10 /* 确保要比 Minimap 组件层级要高，防止被遮挡而点击不到 */,
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
