import { startTransition, useEffect, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  Field,
  type FlowNodeEntity,
  FlowNodeRenderData,
  type NodeRenderReturnType,
  useClientContext,
  useWatchFormState,
} from "@flowgram.ai/fixed-layout-editor";
import { IconCopy, IconDotsVertical, IconExclamationCircle, IconGripVertical, IconLabel, IconX } from "@tabler/icons-react";
import { Button, type ButtonProps, Card, Dropdown, Input, type InputRef, Popover, Tooltip, theme } from "antd";
import { Immer } from "immer";
import { nanoid } from "nanoid";

import { mergeCls } from "@/utils/css";

import { type NodeJSON, type NodeRegistry, NodeType } from "./typings";
import { useNodeRenderContext } from "../NodeRenderContext";

const useInternalRenamingInput = ({ nodeRender }: { nodeRender: NodeRenderReturnType }) => {
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

const InternalNodeCard = ({
  className,
  style,
  children,
  nodeRender,
}: {
  className?: string;
  style?: React.CSSProperties;
  children?: React.ReactNode;
  nodeRender: NodeRenderReturnType;
}) => {
  const nodeRenderData = nodeRender.node.getData(FlowNodeRenderData)!;
  const nodeRegistry = nodeRender.node.getNodeRegistry<NodeRegistry>();

  const isActivated = useMemo(() => nodeRenderData.activated || nodeRenderData.lineActivated, [nodeRenderData.activated, nodeRenderData.lineActivated]);
  const [isHovered, setIsHovered] = useState(false);
  const [isInvalid, setIsInvalid] = useState(false);

  const formState = useWatchFormState(nodeRender.node);
  useEffect(() => setIsInvalid(!!formState?.invalid), [formState?.invalid]);

  return (
    <Card
      className={mergeCls(
        "relative rounded-xl shadow-sm",
        { "border-primary": isActivated },
        nodeRegistry.meta?.clickable ? "cursor-pointer" : "cursor-default",
        className
      )}
      style={style}
      styles={{ body: { padding: 0 } }}
      hoverable
      onMouseEnter={() => startTransition(() => setIsHovered(true))}
      onMouseLeave={() => startTransition(() => setIsHovered(false))}
    >
      <div className="relative z-1">{children}</div>
      <div
        className="absolute z-0 rounded-xl border-solid border-transparent transition-all duration-500"
        style={{
          top: "-1px",
          left: "-1px",
          right: "-1px",
          bottom: "-1px",
          borderWidth: "2px",
          borderColor: isHovered ? "var(--color-primary)" : isInvalid ? "var(--color-error)" : void 0,
        }}
      />
    </Card>
  );
};

const InternalNodeMenuButton = ({
  className,
  style,
  onMenuSelect,
  ...props
}: ButtonProps & {
  onMenuSelect?: (key: "rename" | "duplicate" | "remove") => void;
}) => {
  const { t } = useTranslation();

  const ctx = useClientContext();
  const { operation, playground } = ctx;

  const { node, ...nodeRender } = useNodeRenderContext();
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
    const d1 = node.onEntityChange(callback);
    const d2 = node.parent?.onEntityChange?.(callback);
    return () => {
      d1?.dispose();
      d2?.dispose();
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
            label: nodeRender.isBlockOrderIcon ? t("workflow.detail.design.editor.rename_branch") : t("workflow.detail.design.editor.rename_node"),
            icon: <IconLabel size="1em" />,
            onClick: handleClickRename,
          },
          {
            key: "duplicate",
            label: nodeRender.isBlockOrderIcon ? t("workflow.detail.design.editor.duplicate_branch") : t("workflow.detail.design.editor.duplicate_node"),
            icon: <IconCopy size="1em" />,
            disabled: menuDuplicateDisabled,
            onClick: handleClickDuplicate,
          },
          {
            type: "divider",
          },
          {
            key: "remove",
            label: nodeRender.isBlockOrderIcon ? t("workflow.detail.design.editor.remove_branch") : t("workflow.detail.design.editor.remove_node"),
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

export interface BaseNodeProps {
  className?: string;
  style?: React.CSSProperties;
  children?: React.ReactNode;
  description?: React.ReactNode;
}

export const BaseNode = ({ className, style, children, description }: BaseNodeProps) => {
  const { token: themeToken } = theme.useToken();

  const ctx = useClientContext();
  const { playground } = ctx;

  const nodeRender = useNodeRenderContext();
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

  const [isInvalid, setIsInvalid] = useState(false);

  const formState = useWatchFormState(nodeRender.node);
  useEffect(() => setIsInvalid(!!formState?.invalid), [formState?.invalid]);

  const {
    inputRef,
    visible: inputVisible,
    value: inputValue,
    onClick: handleNodeRenameClick,
    onChange: handleNodeNameChange,
    onConfirm: handleNodeNameConfirm,
  } = useInternalRenamingInput({ nodeRender });

  return (
    <Popover
      classNames={{ root: "shadow-md" }}
      styles={{ body: { padding: 0 } }}
      arrow={false}
      content={
        inputVisible ? null : (
          <InternalNodeMenuButton
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
      placement="rightTop"
    >
      <div className="group/node relative">
        <InternalNodeCard className={mergeCls("w-[320px]", className)} style={style} nodeRender={nodeRender}>
          {children != null ? (
            children
          ) : (
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
                  <Field<string> name="name">{({ field: { value } }) => <>{value || "\u00A0"}</>}</Field>
                </div>
                {description != null && (
                  <div className="truncate text-xs" style={{ color: themeToken.colorTextTertiary }}>
                    {description}
                  </div>
                )}
              </div>
              <div className="flex items-center justify-center" onClick={(e) => e.stopPropagation()}>
                {isInvalid && <IconExclamationCircle color="var(--color-error)" size="1.25em" />}
              </div>
            </div>
          )}
        </InternalNodeCard>

        {!playground.config.readonlyOrDisabled && nodeRegistry.meta?.draggable === true && (
          <div className="absolute top-1/2 -left-4 z-1 hidden -translate-y-1/2 group-hover/node:block">
            <IconGripVertical size="1em" stroke="1" />
          </div>
        )}

        {!playground.config.readonlyOrDisabled && (
          <div className={mergeCls("absolute top-1/2 left-2 right-2 -translate-y-1/2 z-1", inputVisible ? "block" : "hidden")}>
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

export interface BranchNodeProps extends BaseNodeProps {}

export const BranchNode = ({ className, style, children, description }: BranchNodeProps) => {
  const ctx = useClientContext();
  const { playground } = ctx;

  const nodeRender = useNodeRenderContext();
  const nodeRegistry = nodeRender.node.getNodeRegistry<NodeRegistry>();

  const {
    inputRef,
    visible: inputVisible,
    value: inputValue,
    onClick: handleNodeRenameClick,
    onChange: handleNodeNameChange,
    onConfirm: handleNodeNameConfirm,
  } = useInternalRenamingInput({ nodeRender });

  return (
    <Popover
      classNames={{ root: "shadow-md" }}
      styles={{ body: { padding: 0 } }}
      arrow={false}
      content={
        inputVisible ? null : (
          <InternalNodeMenuButton
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
      <div className="group/node relative">
        <InternalNodeCard className={mergeCls("w-[240px]", className)} style={style} nodeRender={nodeRender}>
          {children != null ? (
            children
          ) : (
            <div className={mergeCls("overflow-hidden px-3 py-2", inputVisible ? "invisible" : "visible")}>
              <div className="truncate text-center">
                {description ?? <Field<string> name="name">{({ field: { value } }) => <>{value || "\u00A0"}</>}</Field>}
              </div>
            </div>
          )}
        </InternalNodeCard>

        {!playground.config.readonlyOrDisabled && nodeRegistry.meta?.draggable === true && (
          <div className="absolute top-1/2 -left-4 z-1 hidden -translate-y-1/2 group-hover/node:block">
            <IconGripVertical size="1em" stroke="1" />
          </div>
        )}

        {!playground.config.readonlyOrDisabled && (
          <div className={mergeCls("absolute top-1/2 left-2 right-2 -translate-y-1/2 z-1", inputVisible ? "block" : "hidden")}>
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

// TODO: 应放至领域层
export const duplicateNodeJSON = (node: NodeJSON, options?: { withCopySuffix?: boolean }) => {
  const { produce } = new Immer({ autoFreeze: false });
  const deepClone = (node: NodeJSON, { withCopySuffix, nodeIdMap }: { withCopySuffix: boolean; nodeIdMap: Map<string, string> }) => {
    return produce(node, (draft) => {
      draft.data ??= {};
      draft.id = nanoid();
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
