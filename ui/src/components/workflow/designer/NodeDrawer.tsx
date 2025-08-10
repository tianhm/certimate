import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { type FlowNodeEntity, getNodeForm } from "@flowgram.ai/fixed-layout-editor";
import { IconX } from "@tabler/icons-react";
import { useControllableValue } from "ahooks";
import { Anchor, type AnchorProps, App, Button, Drawer, Flex, Form, Typography } from "antd";
import { isEqual } from "radash";

import { useTriggerElement } from "@/hooks";
import { getErrMsg } from "@/utils/error";

import BizApplyNodeConfigForm from "./forms/BizApplyNodeConfigForm";
import BizDeployNodeConfigForm from "./forms/BizDeployNodeConfigForm";
import BizMonitorNodeConfigForm from "./forms/BizMonitorNodeConfigForm";
import BizNotifyNodeConfigForm from "./forms/BizNotifyNodeConfigForm";
import BizUploadNodeConfigForm from "./forms/BizUploadNodeConfigForm";
import BranchBlockNodeConfigForm from "./forms/BranchBlockNodeConfigForm";
import StartNodeConfigForm from "./forms/StartNodeConfigForm";
import { type NodeRegistry, NodeType } from "./nodes/typings";

export interface NodeDrawerProps {
  children?: React.ReactNode;
  anchor?: Pick<AnchorProps, "items"> | false;
  loading?: boolean;
  node?: FlowNodeEntity;
  open?: boolean;
  trigger?: React.ReactNode;
  onOpenChange?: (open: boolean) => void;
}

const NodeDrawer = (_: NodeDrawerProps) => {
  const { anchor, loading, node, trigger, ...props } = _;

  const { t } = useTranslation();

  const { modal, notification } = App.useApp();

  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const containerRef = useRef<HTMLDivElement>(null);

  const triggerEl = useTriggerElement(trigger, { onClick: () => setOpen(true) });

  const [formInst] = Form.useForm();
  const [formPending, setFormPending] = useState<boolean>(false);

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
  const renderNodeConfigForm = () => {
    if (node == null) return null;

    const formProps = { form: formInst, node: node };
    switch (node.flowNodeType) {
      case NodeType.Start:
        return <StartNodeConfigForm {...formProps} />;
      case NodeType.BizApply:
        return <BizApplyNodeConfigForm {...formProps} />;
      case NodeType.BizUpload:
        return <BizUploadNodeConfigForm {...formProps} />;
      case NodeType.BizMonitor:
        return <BizMonitorNodeConfigForm {...formProps} />;
      case NodeType.BizDeploy:
        return <BizDeployNodeConfigForm {...formProps} />;
      case NodeType.BizNotify:
        return <BizNotifyNodeConfigForm {...formProps} />;
      case NodeType.BranchBlock:
        return <BranchBlockNodeConfigForm {...formProps} />;
      default:
        console.warn(`[certimate] unsupported workflow node type: ${node.flowNodeType}`);
        return null;
    }
  };

  const handleOkClick = async () => {
    if (node == null) {
      setOpen(false);
      return;
    }

    setFormPending(true);
    try {
      await formInst.validateFields();
    } catch (err) {
      setFormPending(false);
      throw err;
    }

    try {
      getNodeForm(node)!.setValueIn("config", formInst.getFieldsValue(true));
      getNodeForm(node)!.validate();

      setOpen(false);
    } catch (err) {
      notification.error({ message: t("common.text.request_error"), description: getErrMsg(err) });

      throw err;
    } finally {
      setFormPending(false);
    }
  };

  const handleCancelClick = () => {
    if (formPending) return;

    setOpen(false);
  };

  const handleClose = () => {
    if (formPending) return;

    const picker = (obj: Record<string, unknown>) => {
      return Object.entries(obj).reduce(
        (acc, [key, value]) => {
          const isEmpty =
            value === null ||
            value === undefined ||
            (typeof value === "string" && value === "") ||
            (Array.isArray(value) && value.length === 0) ||
            (typeof value === "object" && !Array.isArray(value) && Object.keys(value).length === 0);

          if (!isEmpty) {
            acc[key] = value;
          }

          return acc;
        },
        {} as Record<string, unknown>
      );
    };
    const oldValues = picker(node?.toJSON()?.data?.config ?? {});
    const newValues = picker(formInst.getFieldsValue(true));
    const changed = !isEqual(oldValues, {}) && !isEqual(oldValues, newValues);

    const { promise, resolve, reject } = Promise.withResolvers();
    if (changed) {
      console.log(oldValues, newValues);
      modal.confirm({
        title: t("common.text.operation_confirm"),
        content: t("workflow.detail.design.unsaved_changes.confirm"),
        onOk: () => resolve(void 0),
        onCancel: () => reject(),
      });
    } else {
      resolve(void 0);
    }

    promise.then(() => setOpen(false));
  };

  return (
    <>
      {triggerEl}

      <Drawer
        styles={{
          header: {
            paddingBottom: anchor != null && anchor !== false ? 0 : void 0,
          },
        }}
        afterOpenChange={setOpen}
        closeIcon={false}
        destroyOnHidden
        footer={
          <Flex className="px-2" justify="end" gap="small">
            <Button onClick={handleCancelClick}>{t("common.button.cancel")}</Button>
            <Button loading={formPending} type="primary" onClick={handleOkClick}>
              {t("common.button.save")}
            </Button>
          </Flex>
        }
        forceRender={false}
        loading={loading}
        maskClosable={!formPending}
        open={open}
        size="large"
        title={
          <>
            <Flex align="center" justify="space-between" gap="small">
              <div>{renderNodeIcon()}</div>
              <div className="flex-1 truncate">{node?.toJSON()?.data?.name}</div>
              <Button className="ant-drawer-close" style={{ marginInline: 0 }} icon={<IconX size="1.25em" />} size="small" type="text" onClick={handleClose} />
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
        onClose={handleClose}
      >
        <div ref={containerRef}>{renderNodeConfigForm()}</div>
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
      setNode(void 0);
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
