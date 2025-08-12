import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { validDomainName } from "@/utils/validators";

type DeployNodeConfigFormHuaweiCloudOBSConfigFieldValues = Nullish<{
  endpoint: string;
  bucket: string;
  domain: string;
}>;

export interface DeployNodeConfigFormHuaweiCloudOBSConfigProps {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: DeployNodeConfigFormHuaweiCloudOBSConfigFieldValues;
  onValuesChange?: (values: DeployNodeConfigFormHuaweiCloudOBSConfigFieldValues) => void;
}

const initFormModel = (): DeployNodeConfigFormHuaweiCloudOBSConfigFieldValues => {
  return {};
};

const DeployNodeConfigFormHuaweiCloudOBSConfig = ({
  form: formInst,
  formName,
  disabled,
  initialValues,
  onValuesChange,
}: DeployNodeConfigFormHuaweiCloudOBSConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    endpoint: z
      .string(t("workflow_node.deploy.form.huaweicloud_obs_endpoint.placeholder"))
      .nonempty(t("workflow_node.deploy.form.huaweicloud_obs_endpoint.placeholder")),
    bucket: z
      .string(t("workflow_node.deploy.form.huaweicloud_obs_bucket.placeholder"))
      .nonempty(t("workflow_node.deploy.form.huaweicloud_obs_bucket.placeholder")),
    domain: z
      .string(t("workflow_node.deploy.form.huaweicloud_obs_domain.placeholder"))
      .refine((v) => validDomainName(v, { allowWildcard: true }), t("common.errmsg.domain_invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values);
  };

  return (
    <Form
      form={formInst}
      disabled={disabled}
      initialValues={initialValues ?? initFormModel()}
      layout="vertical"
      name={formName}
      onValuesChange={handleFormChange}
    >
      <Form.Item
        name="endpoint"
        label={t("workflow_node.deploy.form.huaweicloud_obs_endpoint.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.huaweicloud_obs_endpoint.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.huaweicloud_obs_endpoint.placeholder")} />
      </Form.Item>

      <Form.Item
        name="bucket"
        label={t("workflow_node.deploy.form.huaweicloud_obs_bucket.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.huaweicloud_obs_bucket.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.huaweicloud_obs_bucket.placeholder")} />
      </Form.Item>

      <Form.Item
        name="domain"
        label={t("workflow_node.deploy.form.huaweicloud_obs_domain.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.huaweicloud_obs_domain.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.huaweicloud_obs_domain.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default DeployNodeConfigFormHuaweiCloudOBSConfig;
