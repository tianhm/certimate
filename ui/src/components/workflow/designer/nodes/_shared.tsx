import { useContext, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useClientContext } from "@flowgram.ai/fixed-layout-editor";
import { IconCopy, IconDotsVertical, IconGripVertical, IconLabel, IconX } from "@tabler/icons-react";
import { Button, type ButtonProps, Card, Dropdown, Popover, Tooltip, Typography, theme } from "antd";
import { mergeCls } from "@/utils/css";

import { type NodeRegistry } from "./typings";
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
  const { playground } = ctx;

  const { node, deleteNode, isBlockIcon, isBlockOrderIcon } = useContext(NodeRenderContext);
  const nodeRegistry = node.getNodeRegistry<NodeRegistry>();

  const getLatestNodeDeleteDisabledState = () => {
    if (nodeRegistry.canDelete != null) {
      return !nodeRegistry.canDelete(ctx, node);
    }
    return !!nodeRegistry.meta?.deleteDisable;
  };
  const [nodeDeleteDisabled, setNodeDeleteDisabled] = useState(() => getLatestNodeDeleteDisabledState());
  useEffect(() => {
    // 这里不能使用 useMemo() 来决定 nodeDeleteDisabled，因为依赖项没有发生改变（对象引用始终是同一个）
    // 因此需要使用 useEffect() 来监听 node 和 node.parent 的变化，并更新 nodeDeleteDisabled 的状态
    const disposable1 = node.onEntityChange(() => setNodeDeleteDisabled(getLatestNodeDeleteDisabledState()));
    const disposable2 = node.parent?.onEntityChange(() => setNodeDeleteDisabled(getLatestNodeDeleteDisabledState()));
    return () => {
      disposable1?.dispose();
      disposable2?.dispose();
    };
  }, []);

  return playground.config.readonlyOrDisabled ? null : (
    <Dropdown
      arrow={false}
      destroyOnHidden
      menu={{
        items: [
          {
            key: "rename",
            label: isBlockIcon || isBlockOrderIcon ? t("workflow.detail.design.nodes.rename_branch") : t("workflow.detail.design.nodes.rename_node"),
            icon: <IconLabel size="1em" />,
            onClick: () => {
              alert("TODO: rename");
            },
          },
          {
            key: "duplicate",
            label: isBlockIcon || isBlockOrderIcon ? t("workflow.detail.design.nodes.duplicate_branch") : t("workflow.detail.design.nodes.duplicate_node"),
            icon: <IconCopy size="1em" />,
            onClick: () => {
              alert("TODO: duplicate");
            },
          },
          {
            type: "divider",
          },
          {
            key: "remove",
            label: isBlockIcon || isBlockOrderIcon ? t("workflow.detail.design.nodes.remove_branch") : t("workflow.detail.design.nodes.remove_node"),
            icon: <IconX size="1em" />,
            danger: true,
            disabled: nodeDeleteDisabled,
            onClick: () => {
              deleteNode();
            },
          },
        ],
      }}
      trigger={["click"]}
    >
      <Button className={className} style={style} icon={<IconDotsVertical color="grey" size="1.25em" />} type="text" {...props} />
    </Dropdown>
  );
};
