import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useClientContext, useNodeRender } from "@flowgram.ai/fixed-layout-editor";
import { IconCopy, IconDotsVertical, IconGripVertical, IconLabel, IconX } from "@tabler/icons-react";
import { Button, type ButtonProps, Card, Dropdown, type MenuProps, Popover, Tooltip, Typography, theme } from "antd";
import { mergeCls } from "@/utils/css";

import { type NodeRegistry } from "./typings";

export const BaseNode = ({ className, style, children }: { className?: string; style?: React.CSSProperties; children?: React.ReactNode }) => {
  const { token: themeToken } = theme.useToken();

  const ctx = useClientContext();
  const { playground } = ctx;

  const nodeRender = useNodeRender();
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

  const nodeRender = useNodeRender();
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
  const ctx = useClientContext();
  const { playground } = ctx;

  const menuItems = useNodeMenuItems();

  return playground.config.readonlyOrDisabled ? null : (
    <Dropdown menu={{ items: menuItems }} trigger={["click"]} arrow={false}>
      <Button className={className} style={style} ghost icon={<IconDotsVertical color="grey" size="1.25em" />} type="text" {...props} />
    </Dropdown>
  );
};

export const useNodeMenuItems = () => {
  const { t } = useTranslation();

  const ctx = useClientContext();
  const { operation } = ctx;

  const nodeRender = useNodeRender();
  const nodeRegistry = nodeRender.node.getNodeRegistry<NodeRegistry>();

  const nodeDeleteDisabled = useMemo(() => {
    if (nodeRegistry.canDelete != null) {
      return !nodeRegistry.canDelete(ctx, nodeRender.node);
    }
    return nodeRegistry.meta!.deleteDisable;
  }, [nodeRegistry, nodeRender.node]);

  const menuItems = useMemo<Required<MenuProps>["items"]>(() => {
    return [
      {
        key: "rename",
        label: t("workflow_node.action.rename_node"),
        icon: <IconLabel size="1em" />,
        onClick: () => {
          operation.deleteNode(nodeRender.node);
          alert("TODO: rename");
        },
      },
      {
        key: "duplicate",
        label: t("workflow_node.action.duplicate_node"),
        icon: <IconCopy size="1em" />,
        onClick: () => {
          operation.deleteNode(nodeRender.node);
          alert("TODO: duplicate");
        },
      },
      {
        type: "divider",
      },
      {
        key: "remove",
        label: t("workflow_node.action.remove_node"),
        icon: <IconX size="1em" />,
        danger: true,
        disabled: nodeDeleteDisabled,
        onClick: () => {
          operation.deleteNode(nodeRender.node);
          alert("TODO: remove");
        },
      },
    ];
  }, [nodeRender.node, nodeDeleteDisabled]);

  return menuItems;
};
