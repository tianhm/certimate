import { startTransition, useCallback, useMemo, useState } from "react";
import { type FlowNodeEntity } from "@flowgram.ai/fixed-layout-editor";
import { useControllableValue } from "ahooks";

import { useTriggerElement } from "@/hooks";

import BizApplyNodeConfigDrawer from "./forms/BizApplyNodeConfigDrawer";
import BizDeployNodeConfigDrawer from "./forms/BizDeployNodeConfigDrawer";
import BizMonitorNodeConfigDrawer from "./forms/BizMonitorNodeConfigDrawer";
import BizNotifyNodeConfigDrawer from "./forms/BizNotifyNodeConfigDrawer";
import BizUploadNodeConfigDrawer from "./forms/BizUploadNodeConfigDrawer";
import BranchBlockNodeConfigDrawer from "./forms/BranchBlockNodeConfigDrawer";
import DelayNodeConfigDrawer from "./forms/DelayNodeConfigDrawer";
import StartNodeConfigDrawer from "./forms/StartNodeConfigDrawer";
import { NodeType } from "./nodes/typings";

export interface NodeDrawerProps {
  afterClose?: () => void;
  children?: React.ReactNode;
  loading?: boolean;
  node?: FlowNodeEntity;
  open?: boolean;
  trigger?: React.ReactNode;
  onOpenChange?: (open: boolean) => void;
}

const NodeDrawer = ({ node, trigger, ...props }: NodeDrawerProps) => {
  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const triggerEl = useTriggerElement(trigger, { onClick: () => setOpen(true) });

  const drawerProps = useMemo(
    () => ({
      ...props,
      node: node!,
      open: open,
      onOpenChange: (open: boolean) => {
        setOpen(open);
      },
    }),
    [props, node, open]
  );

  return (
    <>
      {triggerEl}

      {node?.flowNodeType === NodeType.Start ? (
        <StartNodeConfigDrawer {...drawerProps} />
      ) : node?.flowNodeType === NodeType.Delay ? (
        <DelayNodeConfigDrawer {...drawerProps} />
      ) : node?.flowNodeType === NodeType.BranchBlock ? (
        <BranchBlockNodeConfigDrawer {...drawerProps} />
      ) : node?.flowNodeType === NodeType.BizApply ? (
        <BizApplyNodeConfigDrawer {...drawerProps} />
      ) : node?.flowNodeType === NodeType.BizUpload ? (
        <BizUploadNodeConfigDrawer {...drawerProps} />
      ) : node?.flowNodeType === NodeType.BizMonitor ? (
        <BizMonitorNodeConfigDrawer {...drawerProps} />
      ) : node?.flowNodeType === NodeType.BizDeploy ? (
        <BizDeployNodeConfigDrawer {...drawerProps} />
      ) : node?.flowNodeType === NodeType.BizNotify ? (
        <BizNotifyNodeConfigDrawer {...drawerProps} />
      ) : (
        <></>
      )}
    </>
  );
};

const useDrawer = () => {
  type NodeDataType = NodeDrawerProps["node"];
  const [node, setNode] = useState<NodeDataType>();
  const [open, setOpen] = useState<boolean>(false);

  const onOpenChange = useCallback((open: boolean) => {
    setOpen(open);
  }, []);

  return {
    drawerProps: {
      afterClose: () => {
        startTransition(() => {
          if (!open) {
            setNode(void 0);
          }
        });
      },
      node,
      open,
      onOpenChange,
    },

    open: (node: NonNullable<NodeDataType>) => {
      setNode(node);
      setOpen(true);
    },
    close: () => {
      setOpen(false);
    },
  };
};

const _default = Object.assign(NodeDrawer, {
  useDrawer,
});

export default _default;
