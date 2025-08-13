import { useTranslation } from "react-i18next";
import { type FlowNodeEntity } from "@flowgram.ai/fixed-layout-editor";
import { Form } from "antd";

import { NodeConfigDrawer } from "./_shared";
import BizNotifyNodeConfigForm from "./BizNotifyNodeConfigForm";
import { NodeType } from "../nodes/typings";

export interface BizNotifyNodeConfigDrawerProps {
  loading?: boolean;
  node: FlowNodeEntity;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}

const BizNotifyNodeConfigDrawer = ({ node, ...props }: BizNotifyNodeConfigDrawerProps) => {
  if (node.flowNodeType !== NodeType.BizNotify) {
    console.warn(`[certimate] current workflow node type is not: ${NodeType.BizNotify}`);
  }

  const { i18n } = useTranslation();

  const [formInst] = Form.useForm();

  return (
    <NodeConfigDrawer
      anchor={{
        items: BizNotifyNodeConfigForm.getAnchorItems({ i18n }),
      }}
      form={formInst}
      node={node}
      {...props}
    >
      <BizNotifyNodeConfigForm form={formInst} node={node} />
    </NodeConfigDrawer>
  );
};

export default BizNotifyNodeConfigDrawer;
