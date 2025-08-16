import { useMemo } from "react";
import { getI18n, useTranslation } from "react-i18next";
import { type FlowNodeEntity, getNodeForm } from "@flowgram.ai/fixed-layout-editor";
import { type AnchorProps, Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { validateCertificate, validatePrivateKey } from "@/api/certificates";
import TextFileInput from "@/components/TextFileInput";
import { type WorkflowNodeConfigForBizUpload, defaultNodeConfigForBizUpload } from "@/domain/workflow";
import { useAntdForm } from "@/hooks";
import { getErrMsg } from "@/utils/error";

import { NodeFormContextProvider } from "./_context";
import { NodeType } from "../nodes/typings";

export interface BizUploadNodeConfigFormProps {
  form: FormInstance;
  node: FlowNodeEntity;
}

const BizUploadNodeConfigForm = ({ node, ...props }: BizUploadNodeConfigFormProps) => {
  if (node.flowNodeType !== NodeType.BizUpload) {
    console.warn(`[certimate] current workflow node type is not: ${NodeType.BizUpload}`);
  }

  const { i18n, t } = useTranslation();

  const initialValues = useMemo(() => {
    return getNodeForm(node)?.getValueIn("config") as WorkflowNodeConfigForBizUpload | undefined;
  }, [node]);

  const formSchema = getSchema({ i18n });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    form: props.form,
    name: "workflowNodeBizUploadConfigForm",
    initialValues: initialValues ?? getInitialValues(),
  });

  const handleCertificateChange = async (value: string) => {
    try {
      const resp = await validateCertificate(value);
      formInst.setFields([
        {
          name: "domains",
          value: resp.data.domains,
        },
        {
          name: "certificate",
          value: value,
        },
      ]);
    } catch (e) {
      formInst.setFields([
        {
          name: "domains",
          value: "",
        },
        {
          name: "certificate",
          value: value,
          errors: [getErrMsg(e)],
        },
      ]);
    }
  };

  const handlePrivateKeyChange = async (value: string) => {
    try {
      await validatePrivateKey(value);
      formInst.setFields([
        {
          name: "privateKey",
          value: value,
        },
      ]);
    } catch (e) {
      formInst.setFields([
        {
          name: "privateKey",
          value: value,
          errors: [getErrMsg(e)],
        },
      ]);
    }
  };

  return (
    <NodeFormContextProvider value={{ node }}>
      <Form {...formProps} clearOnDestroy={true} form={formInst} layout="vertical" preserve={false} scrollToFirstError>
        <div id="parameters" data-anchor="parameters">
          <Form.Item name="domains" label={t("workflow_node.upload.form.domains.label")} rules={[formRule]}>
            <Input variant="filled" placeholder={t("workflow_node.upload.form.domains.placeholder")} readOnly />
          </Form.Item>

          <Form.Item name="certificate" label={t("workflow_node.upload.form.certificate.label")} rules={[formRule]}>
            <TextFileInput
              autoSize={{ minRows: 3, maxRows: 10 }}
              placeholder={t("workflow_node.upload.form.certificate.placeholder")}
              onChange={handleCertificateChange}
            />
          </Form.Item>

          <Form.Item name="privateKey" label={t("workflow_node.upload.form.private_key.label")} rules={[formRule]}>
            <TextFileInput
              autoSize={{ minRows: 3, maxRows: 10 }}
              placeholder={t("workflow_node.upload.form.private_key.placeholder")}
              onChange={handlePrivateKeyChange}
            />
          </Form.Item>
        </div>
      </Form>
    </NodeFormContextProvider>
  );
};

const getAnchorItems = ({ i18n = getI18n() }: { i18n?: ReturnType<typeof getI18n> }): Required<AnchorProps>["items"] => {
  const { t } = i18n;

  return ["parameters"].map((key) => ({
    key: key,
    title: t(`workflow_node.upload.form_anchor.${key}.tab`),
    href: "#" + key,
  }));
};

const getInitialValues = (): Nullish<z.infer<ReturnType<typeof getSchema>>> => {
  return {
    certificate: "",
    privateKey: "",
    ...defaultNodeConfigForBizUpload(),
  };
};

const getSchema = ({ i18n = getI18n() }: { i18n?: ReturnType<typeof getI18n> }) => {
  const { t } = i18n;

  return z.object({
    domains: z.string().nullish(),
    certificate: z
      .string()
      .min(1, t("workflow_node.upload.form.certificate.placeholder"))
      .max(20480, t("common.errmsg.string_max", { max: 20480 })),
    privateKey: z
      .string()
      .min(1, t("workflow_node.upload.form.private_key.placeholder"))
      .max(20480, t("common.errmsg.string_max", { max: 20480 })),
  });
};

const _default = Object.assign(BizUploadNodeConfigForm, {
  getAnchorItems,
  getSchema,
});

export default _default;
