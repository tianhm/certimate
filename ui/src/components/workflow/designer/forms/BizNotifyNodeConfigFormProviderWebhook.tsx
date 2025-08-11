import { getI18n, useTranslation } from "react-i18next";
import { Form } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import CodeInput from "@/components/CodeInput";
import Tips from "@/components/Tips";

import { useFormNestedFieldsContext } from "./_context";

const BizNotifyNodeConfigFormProviderWebhook = () => {
  const { i18n, t } = useTranslation();

  const { parentNamePath } = useFormNestedFieldsContext();
  const formSchema = z.object({
    [parentNamePath]: getSchema({ i18n }),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const formInst = Form.useFormInstance();
  const initialValues = getInitialValues();

  const handleWebhookDataBlur = () => {
    const value = formInst.getFieldValue("webhookData");
    try {
      const json = JSON.stringify(JSON.parse(value), null, 2);
      formInst.setFieldValue("webhookData", json);
    } catch {
      return;
    }
  };

  return (
    <>
      <Form.Item
        name={[parentNamePath, "webhookData"]}
        initialValue={initialValues.webhookData}
        label={t("workflow_node.notify.form.webhook_data.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.notify.form.webhook_data.tooltip") }}></span>}
      >
        <CodeInput
          height="auto"
          minHeight="64px"
          maxHeight="256px"
          language="json"
          placeholder={t("workflow_node.notify.form.webhook_data.placeholder")}
          onBlur={handleWebhookDataBlur}
        />
      </Form.Item>

      <Form.Item>
        <Tips message={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.notify.form.webhook_data.guide") }}></span>} />
      </Form.Item>
    </>
  );
};

const getInitialValues = (): Nullish<z.infer<ReturnType<typeof getSchema>>> => {
  return {};
};

const getSchema = ({ i18n = getI18n() }: { i18n?: ReturnType<typeof getI18n> }) => {
  const { t } = i18n;

  return z.object({
    webhookData: z
      .string()
      .nullish()
      .refine((v) => {
        if (!v) return true;

        try {
          const obj = JSON.parse(v);
          return typeof obj === "object" && !Array.isArray(obj);
        } catch {
          return false;
        }
      }, t("workflow_node.notify.form.webhook_data.errmsg.json_invalid")),
  });
};

const _default = Object.assign(BizNotifyNodeConfigFormProviderWebhook, {
  getInitialValues,
  getSchema,
});

export default _default;
