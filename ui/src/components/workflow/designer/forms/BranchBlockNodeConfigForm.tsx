import { useMemo, useRef } from "react";
import { getI18n, useTranslation } from "react-i18next";
import { type FlowNodeEntity, getNodeForm } from "@flowgram.ai/fixed-layout-editor";
import { type AnchorProps, Form, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type WorkflowNodeConfigForCondition, defaultNodeConfigForCondition } from "@/domain/workflow";
import { useAntdForm } from "@/hooks";

import { NodeFormContextProvider } from "./_context";
import BranchBlockNodeConfigFormExpressionEditor, { type BranchBlockNodeConfigFormExpressionEditorInstance } from "./BranchBlockNodeConfigFormExpressionEditor";

import { NodeType } from "../nodes/typings";

export interface BranchBlockNodeConfigFormProps {
  form: FormInstance;
  node: FlowNodeEntity;
}

const BranchBlockNodeConfigForm = ({ node, ...props }: BranchBlockNodeConfigFormProps) => {
  if (node.flowNodeType !== NodeType.BranchBlock) {
    console.warn(`[certimate] current workflow node type is not: ${NodeType.BranchBlock}`);
  }

  const { i18n, t } = useTranslation();

  const initialValues = useMemo(() => {
    return getNodeForm(node)?.getValueIn("config") as WorkflowNodeConfigForCondition | undefined;
  }, [node]);

  const formSchema = getSchema({ i18n }).superRefine(async (values, ctx) => {
    if (values.expression != null) {
      try {
        await exprEditorRef.current!.validate();
      } catch {
        ctx.addIssue({
          code: "custom",
          message: t("workflow_node.branch_block.form.expression.errmsg.invalid"),
          path: ["expression"],
        });
      }
    }
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    form: props.form,
    name: "workflowNodeBranchBlockConfigForm",
    initialValues: initialValues ?? getInitialValues(),
  });

  const exprEditorRef = useRef<BranchBlockNodeConfigFormExpressionEditorInstance>(null);

  return (
    <NodeFormContextProvider value={{ node }}>
      <Form {...formProps} clearOnDestroy={true} form={formInst} layout="vertical" preserve={false} scrollToFirstError>
        <div id="parameters" data-anchor="parameters">
          <Form.Item name="expression" label={t("workflow_node.condition.form.expression.label")} rules={[formRule]}>
            <BranchBlockNodeConfigFormExpressionEditor ref={exprEditorRef} />
          </Form.Item>
        </div>
      </Form>
    </NodeFormContextProvider>
  );
};

const getAnchorItems = ({ i18n = getI18n() }: { i18n: ReturnType<typeof getI18n> }): Required<AnchorProps>["items"] => {
  const { t } = i18n;

  return ["parameters"].map((key) => ({
    key: key,
    title: t(`workflow_node.branch_block.form_anchor.${key}.tab`),
    href: "#" + key,
  }));
};

const getInitialValues = (): Nullish<z.infer<ReturnType<typeof getSchema>>> => {
  return defaultNodeConfigForCondition();
};

const getSchema = ({ i18n = getI18n() }: { i18n: ReturnType<typeof getI18n> }) => {
  const { t: _ } = i18n;

  return z.object({
    expression: z.any().nullish(),
  });
};

const _default = Object.assign(BranchBlockNodeConfigForm, {
  getAnchorItems,
  getSchema,
});

export default _default;
