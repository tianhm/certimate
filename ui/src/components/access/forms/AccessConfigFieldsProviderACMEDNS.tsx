import { getI18n, useTranslation } from "react-i18next";
import { Form, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod/v4";

import { useFormNestedFieldsContext } from "./_context";

const AccessConfigFieldsProviderACMEDNS = () => {
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
        name="apiBase"
        initialValue={initialValues.apiBase}
        label={t("access.form.acmedns_api_base.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.acmedns_api_base.tooltip") }}></span>}
      >
        <Input placeholder={t("access.form.acmedns_api_base.placeholder")} />
      </Form.Item>

      <Form.Item
        name="storageBaseUrl"
        initialValue={initialValues.storageBaseUrl}
        label={t("access.form.acmedns_storage_base_url.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.acmedns_storage_base_url.tooltip") }}></span>}
      >
        <Input allowClear placeholder={t("access.form.acmedns_storage_base_url.placeholder")} />
      </Form.Item>

      <Form.Item
        name="storagePath"
        initialValue={initialValues.storagePath}
        label={t("access.form.acmedns_storage_path.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.acmedns_storage_path.tooltip") }}></span>}
      >
        <Input allowClear placeholder={t("access.form.acmedns_storage_path.placeholder")} />
      </Form.Item>
    </>
  );
};

const getInitialValues = (): Nullish<z.infer<ReturnType<typeof getSchema>>> => {
  return {
    apiBase: "https://auth.acme-dns.io/",
    storageBaseUrl: "",
    storagePath: "",
  };
};

const getSchema = ({ i18n = getI18n() }: { i18n: ReturnType<typeof getI18n> }) => {
  const { t } = i18n;

  return z.object({
    apiBase: z.url(t("common.errmsg.url_invalid")),
    storageBaseUrl: z
      .string()
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .nullish(),
    storagePath: z
      .string()
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .nullish(),
  });
};

const _default = Object.assign(AccessConfigFieldsProviderACMEDNS, {
  getInitialValues,
  getSchema,
});

export default _default;
