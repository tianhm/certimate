import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { type FlowNodeEntity } from "@flowgram.ai/fixed-layout-editor";
import { IconX } from "@tabler/icons-react";
import { useControllableValue } from "ahooks";
import { Anchor, type AnchorProps, App, Button, Drawer, Flex, Typography } from "antd";

import { useTriggerElement } from "@/hooks";
import { getErrMsg } from "@/utils/error";

import { type NodeRegistry } from "./nodes/typings";

export interface NodeDrawerProps {
  children?: React.ReactNode;
  anchor?: Pick<AnchorProps, "items"> | false;
  confirmLoading?: boolean;
  loading?: boolean;
  footer?: boolean;
  node?: FlowNodeEntity;
  open?: boolean;
  trigger?: React.ReactNode;
  onOpenChange?: (open: boolean) => void;
}

const NodeDrawer = (_: NodeDrawerProps) => {
  const { children, anchor, confirmLoading, footer = true, loading, node, trigger, ...props } = _;

  const { t } = useTranslation();

  const { notification } = App.useApp();

  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const containerRef = useRef<HTMLDivElement>(null);

  const triggerEl = useTriggerElement(trigger, { onClick: () => setOpen(true) });

  // const nodeRender = useNodeRender(node);
  const nodeRegistry = node?.getNodeRegistry<NodeRegistry>();
  const NodeIcon = nodeRegistry?.meta?.icon;
  const renderNodeIcon = () =>
    NodeIcon == null ? null : (
      <div
        className="flex size-6 items-center justify-center rounded-lg bg-white text-primary shadow-md dark:bg-stone-200"
        style={{
          color: nodeRegistry?.meta?.iconColor,
          backgroundColor: nodeRegistry?.meta?.iconBgColor,
        }}
      >
        <NodeIcon size="1em" color={nodeRegistry?.meta?.iconColor} stroke="1.25" />
      </div>
    );

  const handleOkClick = async () => {
    setFormPending(true);
    try {
      // TODO:
    } catch (err) {
      setFormPending(false);
      throw err;
    }

    try {
      // TODO:

      setOpen(false);
    } catch (err) {
      notification.error({ message: t("common.text.request_error"), description: getErrMsg(err) });

      throw err;
    } finally {
      setFormPending(false);
    }
  };

  const handleCancelClick = () => {
    if (confirmLoading) return;

    setOpen(false);
  };

  return (
    <>
      {triggerEl}

      <Drawer
        styles={{
          header: {
            paddingBottom: anchor != null && anchor !== false ? 0 : undefined,
          },
        }}
        afterOpenChange={setOpen}
        closeIcon={false}
        destroyOnHidden
        footer={
          footer ? (
            <Flex className="px-2" justify="end" gap="small">
              <Button onClick={handleCancelClick}>{t("common.button.cancel")}</Button>
              <Button loading={confirmLoading} type="primary" onClick={handleOkClick}>
                {t("common.button.save")}
              </Button>
            </Flex>
          ) : (
            false
          )
        }
        loading={loading}
        maskClosable={!confirmLoading}
        open={open}
        size="large"
        title={
          <>
            <Flex align="center" justify="space-between" gap="small">
              <div>{renderNodeIcon()}</div>
              <div className="flex-1 truncate">{node?.toJSON()?.data?.name}</div>
              <Button
                className="ant-drawer-close"
                style={{ marginInline: 0 }}
                icon={<IconX size="1.25em" />}
                size="small"
                type="text"
                onClick={handleCancelClick}
              />
            </Flex>
            <div className="mt-3 text-sm font-normal">
              <Typography.Text className="text-xs" type="secondary">
                <span>{t("workflow.detail.design.drawer.node_id.label")}</span>
                <span>{node?.id}</span>
              </Typography.Text>
            </div>
            {anchor != null && anchor !== false && (
              <div className="mt-3 text-sm font-normal">
                <Anchor affix={false} getContainer={() => containerRef.current!} direction="horizontal" items={anchor.items} />
              </div>
            )}
          </>
        }
        onClose={handleCancelClick}
      >
        <div ref={containerRef}>{children}</div>
      </Drawer>
    </>
  );
};

const useProps = () => {
  const [node, setNode] = useState<NodeDrawerProps["node"]>();
  const [open, setOpen] = useState<boolean>(false);

  const onOpenChange = (open: boolean) => {
    setOpen(open);

    if (!open) {
      setNode(undefined);
    }
  };

  return {
    node,
    open,
    setNode,
    setOpen,
    onOpenChange,
  };
};

const _default = Object.assign(NodeDrawer, {
  useProps,
});

export default _default;
