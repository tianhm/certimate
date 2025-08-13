import { getI18n, useTranslation } from "react-i18next";
import { Form, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { useFormNestedFieldsContext } from "./_context";

const BizDeployNodeConfigFormProviderKubernetesSecret = () => {
  const { i18n, t } = useTranslation();

  const { parentNamePath } = useFormNestedFieldsContext();
  const formSchema = z.object({
    [parentNamePath]: getSchema({ i18n }),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const initialValues = getInitialValues();

  return (
    <>
      <Form.Item
        name={[parentNamePath, "namespace"]}
        initialValue={initialValues.namespace}
        label={t("workflow_node.deploy.form.k8s_namespace.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.k8s_namespace.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.k8s_namespace.placeholder")} />
      </Form.Item>

      <Form.Item
        name={[parentNamePath, "secretName"]}
        initialValue={initialValues.secretName}
        label={t("workflow_node.deploy.form.k8s_secret_name.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.k8s_secret_name.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.k8s_secret_name.placeholder")} />
      </Form.Item>

      <Form.Item
        name={[parentNamePath, "secretType"]}
        initialValue={initialValues.secretType}
        label={t("workflow_node.deploy.form.k8s_secret_type.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.k8s_secret_type.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.k8s_secret_type.placeholder")} />
      </Form.Item>

      <Form.Item
        name={[parentNamePath, "secretDataKeyForCrt"]}
        initialValue={initialValues.secretDataKeyForCrt}
        label={t("workflow_node.deploy.form.k8s_secret_data_key_for_crt.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.k8s_secret_data_key_for_crt.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.k8s_secret_data_key_for_crt.placeholder")} />
      </Form.Item>

      <Form.Item
        name={[parentNamePath, "secretDataKeyForKey"]}
        initialValue={initialValues.secretDataKeyForKey}
        label={t("workflow_node.deploy.form.k8s_secret_data_key_for_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.k8s_secret_data_key_for_key.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.k8s_secret_data_key_for_key.placeholder")} />
      </Form.Item>
    </>
  );
};

const getInitialValues = (): Nullish<z.infer<ReturnType<typeof getSchema>>> => {
  return {
    namespace: "default",
    secretType: "kubernetes.io/tls",
    secretDataKeyForCrt: "tls.crt",
    secretDataKeyForKey: "tls.key",
  };
};

const getSchema = ({ i18n = getI18n() }: { i18n?: ReturnType<typeof getI18n> }) => {
  const { t } = i18n;

  return z.object({
    namespace: z.string().nonempty(t("workflow_node.deploy.form.k8s_namespace.placeholder")),
    secretName: z.string().nonempty(t("workflow_node.deploy.form.k8s_secret_name.placeholder")),
    secretType: z.string().nonempty(t("workflow_node.deploy.form.k8s_secret_type.placeholder")),
    secretDataKeyForCrt: z.string().nonempty(t("workflow_node.deploy.form.k8s_secret_data_key_for_crt.placeholder")),
    secretDataKeyForKey: z.string().nonempty(t("workflow_node.deploy.form.k8s_secret_data_key_for_key.placeholder")),
  });
};

const _default = Object.assign(BizDeployNodeConfigFormProviderKubernetesSecret, {
  getInitialValues,
  getSchema,
});

export default _default;
