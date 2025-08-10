import { memo, useEffect, useMemo, useState } from "react";
import { getI18n, useTranslation } from "react-i18next";
import { Link } from "react-router";
import { QuestionCircleOutlined as IconQuestionCircleOutlined } from "@ant-design/icons";
import { type FlowNodeEntity, getNodeForm } from "@flowgram.ai/fixed-layout-editor";
import { IconChevronRight, IconCircleMinus, IconPlus } from "@tabler/icons-react";
import { useControllableValue } from "ahooks";
import { AutoComplete, Button, Divider, Flex, Form, type FormInstance, Input, InputNumber, Select, Switch, Tooltip, Typography } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import AccessEditDrawer from "@/components/access/AccessEditDrawer";
import AccessSelect from "@/components/access/AccessSelect";
import MultipleSplitValueInput from "@/components/MultipleSplitValueInput";
import ACMEDns01ProviderSelect from "@/components/provider/ACMEDns01ProviderSelect";
import CAProviderSelect from "@/components/provider/CAProviderSelect";
import Show from "@/components/Show";
import { ACCESS_USAGES, ACME_DNS01_PROVIDERS, accessProvidersMap, acmeDns01ProvidersMap, caProvidersMap } from "@/domain/provider";
import { type WorkflowNodeConfigForApply, defaultNodeConfigForApply } from "@/domain/workflow";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { useAccessesStore } from "@/stores/access";
import { useContactEmailsStore } from "@/stores/contact";
import { validDomainName, validIPv4Address, validIPv6Address } from "@/utils/validators";

import { FormNestedFieldsContextProvider, NodeFormContextProvider } from "./_context";
import BizApplyNodeConfigFormProviderAliyunESA from "./BizApplyNodeConfigFormProviderAliyunESA";
import BizApplyNodeConfigFormProviderAWSRoute53 from "./BizApplyNodeConfigFormProviderAWSRoute53";
import BizApplyNodeConfigFormProviderHuaweiCloudDNS from "./BizApplyNodeConfigFormProviderHuaweiCloudDNS";
import BizApplyNodeConfigFormProviderJDCloudDNS from "./BizApplyNodeConfigFormProviderJDCloudDNS";
import BizApplyNodeConfigFormProviderTencentCloudEO from "./BizApplyNodeConfigFormProviderTencentCloudEO";
import { NodeType } from "../nodes/typings";

const MULTIPLE_INPUT_SEPARATOR = ";";

export interface BizApplyNodeConfigFormProps {
  form: FormInstance;
  node: FlowNodeEntity;
}

const BizApplyNodeConfigForm = ({ node, ...props }: BizApplyNodeConfigFormProps) => {
  if (node.flowNodeType !== NodeType.BizApply) {
    console.warn(`[certimate] current workflow node type is not: ${NodeType.BizApply}`);
  }

  const { i18n, t } = useTranslation();

  const { accesses } = useAccessesStore(useZustandShallowSelector("accesses"));

  const initialValues = useMemo(() => {
    return getNodeForm(node)?.getValueIn("config") as WorkflowNodeConfigForApply | undefined;
  }, [node]);

  const formSchema = getSchema({ i18n });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    form: props.form,
    name: "workflowNodeBizApplyConfigForm",
    initialValues: initialValues ?? getInitialValues(),
  });

  const fieldProvider = Form.useWatch<string>("provider", { form: formInst, preserve: true });
  const fieldProviderAccessId = Form.useWatch<string>("providerAccessId", { form: formInst, preserve: true });
  const fieldCAProvider = Form.useWatch<string>("caProvider", { form: formInst, preserve: true });

  const NestedProviderConfigFields = useMemo(() => {
    /*
      注意：如果追加新的子组件，请保持以 ASCII 排序。
      NOTICE: If you add new child component, please keep ASCII order.
      */
    switch (fieldProvider) {
      case ACME_DNS01_PROVIDERS.ALIYUN_ESA:
        return BizApplyNodeConfigFormProviderAliyunESA;
      case ACME_DNS01_PROVIDERS.AWS:
      case ACME_DNS01_PROVIDERS.AWS_ROUTE53:
        return BizApplyNodeConfigFormProviderAWSRoute53;
      case ACME_DNS01_PROVIDERS.HUAWEICLOUD:
      case ACME_DNS01_PROVIDERS.HUAWEICLOUD_DNS:
        return BizApplyNodeConfigFormProviderHuaweiCloudDNS;
      case ACME_DNS01_PROVIDERS.JDCLOUD:
      case ACME_DNS01_PROVIDERS.JDCLOUD_DNS:
        return BizApplyNodeConfigFormProviderJDCloudDNS;
      case ACME_DNS01_PROVIDERS.TENCENTCLOUD_EO:
        return BizApplyNodeConfigFormProviderTencentCloudEO;
    }
  }, [fieldProvider]);

  const [showProvider, setShowProvider] = useState(false);
  useEffect(() => {
    // 通常情况下每个授权信息只对应一个 DNS 提供商，此时无需显示 DNS 提供商字段；
    // 如果对应多个（如 AWS 的 Route53、Lightsail，阿里云的 DNS、ESA，腾讯云的 DNS、EdgeOne 等），则显示。
    if (fieldProviderAccessId) {
      const access = accesses.find((e) => e.id === fieldProviderAccessId);
      const providers = Array.from(acmeDns01ProvidersMap.values()).filter((e) => e.provider === access?.provider);
      setShowProvider(providers.length > 1);
    } else {
      setShowProvider(false);
    }
  }, [accesses, fieldProviderAccessId]);

  const [showCAProviderAccess, setShowCAProviderAccess] = useState(false);
  useEffect(() => {
    // 内置的 CA 提供商（如 Let's Encrypt）无需显示授权信息字段
    if (fieldCAProvider) {
      const provider = caProvidersMap.get(fieldCAProvider);
      setShowCAProviderAccess(!provider?.builtin);
    } else {
      setShowCAProviderAccess(false);
    }
  }, [fieldCAProvider]);

  const handleProviderSelect = (value: string) => {
    if (fieldProvider === value) return;

    // 切换 DNS 提供商时联动授权信息
    if (initialValues?.provider === value) {
      formInst.setFieldValue("providerAccessId", initialValues?.providerAccessId);
    } else {
      if (acmeDns01ProvidersMap.get(fieldProvider)?.provider !== acmeDns01ProvidersMap.get(value)?.provider) {
        formInst.setFieldValue("providerAccessId", void 0);
      }
    }
  };

  const handleProviderAccessSelect = (value: string) => {
    // 切换授权信息时联动 DNS 提供商
    const access = accesses.find((access) => access.id === value);
    const provider = Array.from(acmeDns01ProvidersMap.values()).find((provider) => provider.provider === access?.provider);
    if (fieldProvider !== provider?.type) {
      formInst.setFieldValue("provider", provider?.type);
    }
  };

  const handleCAProviderSelect = (value?: string | undefined) => {
    // 切换 CA 提供商时联动授权信息
    if (value === "") {
      setTimeout(() => {
        formInst.setFieldValue("caProvider", void 0);
        formInst.setFieldValue("caProviderAccessId", void 0);
      }, 1);
    } else if (initialValues?.caProvider === value) {
      formInst.setFieldValue("caProviderAccessId", initialValues?.caProviderAccessId);
    } else {
      if (caProvidersMap.get(fieldCAProvider)?.provider !== caProvidersMap.get(value!)?.provider) {
        formInst.setFieldValue("caProviderAccessId", void 0);
      }
    }
  };

  return (
    <NodeFormContextProvider value={{ node }}>
      <Form {...formProps} clearOnDestroy={true} form={formInst} layout="vertical" preserve={false} scrollToFirstError>
        <Form.Item
          name="domains"
          label={t("workflow_node.apply.form.domains.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.domains.tooltip") }}></span>}
        >
          <MultipleSplitValueInput
            modalTitle={t("workflow_node.apply.form.domains.multiple_input_modal.title")}
            placeholder={t("workflow_node.apply.form.domains.placeholder")}
            placeholderInModal={t("workflow_node.apply.form.domains.multiple_input_modal.placeholder")}
            separator={MULTIPLE_INPUT_SEPARATOR}
            splitOptions={{ removeEmpty: true, trimSpace: true }}
          />
        </Form.Item>

        <Form.Item
          name="contactEmail"
          label={t("workflow_node.apply.form.contact_email.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.contact_email.tooltip") }}></span>}
        >
          <EmailInput placeholder={t("workflow_node.apply.form.contact_email.placeholder")} />
        </Form.Item>

        <Form.Item name="challengeType" label={t("workflow_node.apply.form.challenge_type.label")} rules={[formRule]} hidden>
          <Select
            options={["DNS-01"].map((e) => ({
              label: e,
              value: e.toLowerCase(),
            }))}
            placeholder={t("workflow_node.apply.form.challenge_type.placeholder")}
          />
        </Form.Item>

        <Form.Item name="provider" label={t("workflow_node.apply.form.provider.label")} hidden={!showProvider} rules={[formRule]}>
          <ACMEDns01ProviderSelect
            disabled={!showProvider}
            placeholder={t("workflow_node.apply.form.provider.placeholder")}
            showSearch
            onFilter={(_, option) => {
              if (fieldProviderAccessId) {
                return accesses.find((e) => e.id === fieldProviderAccessId)?.provider === option.provider;
              }

              return true;
            }}
            onSelect={handleProviderSelect}
          />
        </Form.Item>

        <Form.Item noStyle>
          <label className="mb-1 block">
            <div className="flex w-full items-center justify-between gap-4">
              <div className="max-w-full grow truncate">
                <span>{t("workflow_node.apply.form.provider_access.label")}</span>
                <Tooltip title={t("workflow_node.apply.form.provider_access.tooltip")}>
                  <Typography.Text className="ms-1" type="secondary">
                    <IconQuestionCircleOutlined />
                  </Typography.Text>
                </Tooltip>
              </div>
              <div className="text-right">
                <AccessEditDrawer
                  mode="create"
                  trigger={
                    <Button size="small" type="link">
                      {t("workflow_node.apply.form.provider_access.button")}
                      <IconPlus size="1.25em" />
                    </Button>
                  }
                  usage="dns"
                  afterSubmit={(record) => {
                    const provider = accessProvidersMap.get(record.provider);
                    if (provider?.usages?.includes(ACCESS_USAGES.DNS)) {
                      formInst.setFieldValue("providerAccessId", record.id);
                      handleProviderAccessSelect(record.id);
                    }
                  }}
                />
              </div>
            </div>
          </label>
          <Form.Item name="providerAccessId" rules={[formRule]}>
            <AccessSelect
              placeholder={t("workflow_node.apply.form.provider_access.placeholder")}
              showSearch
              onChange={handleProviderAccessSelect}
              onFilter={(_, option) => {
                if (option.reserve) return false;

                const provider = accessProvidersMap.get(option.provider);
                return !!provider?.usages?.includes(ACCESS_USAGES.DNS);
              }}
            />
          </Form.Item>
        </Form.Item>

        <FormNestedFieldsContextProvider value={{ parentNamePath: "providerConfig" }}>
          {NestedProviderConfigFields && <NestedProviderConfigFields />}
        </FormNestedFieldsContextProvider>

        <Divider size="small">
          <Typography.Text className="text-xs font-normal" type="secondary">
            {t("workflow_node.apply.form.certificate_config.label")}
          </Typography.Text>
        </Divider>

        <Form.Item noStyle>
          <label className="mb-1 block">
            <div className="flex w-full items-center justify-between gap-4">
              <div className="max-w-full grow truncate">
                <span>{t("workflow_node.apply.form.ca_provider.label")}</span>
              </div>
              <div className="text-right">
                <Show when={!fieldCAProvider}>
                  <Link className="ant-typography" to="/settings/ssl-provider" target="_blank">
                    <Button size="small" type="link">
                      {t("workflow_node.apply.form.ca_provider.button")}
                      <IconChevronRight size="1.25em" />
                    </Button>
                  </Link>
                </Show>
              </div>
            </div>
          </label>
          <Form.Item name="caProvider" rules={[formRule]}>
            <CAProviderSelect
              allowClear
              placeholder={t("workflow_node.apply.form.ca_provider.placeholder")}
              showSearch
              onSelect={handleCAProviderSelect}
              onClear={handleCAProviderSelect}
            />
          </Form.Item>
        </Form.Item>

        <Form.Item hidden={!showCAProviderAccess} noStyle>
          <label className="mb-1 block">
            <div className="flex w-full items-center justify-between gap-4">
              <div className="max-w-full grow truncate">
                <span>{t("workflow_node.apply.form.ca_provider_access.label")}</span>
              </div>
              <div className="text-right">
                <AccessEditDrawer
                  data={{ provider: caProvidersMap.get(fieldCAProvider!)?.provider }}
                  mode="create"
                  trigger={
                    <Button size="small" type="link">
                      {t("workflow_node.apply.form.ca_provider_access.button")}
                      <IconChevronRight size="1.25em" />
                    </Button>
                  }
                  usage="ca"
                  afterSubmit={(record) => {
                    const provider = accessProvidersMap.get(record.provider);
                    if (provider?.usages?.includes(ACCESS_USAGES.CA)) {
                      formInst.setFieldValue("caProviderAccessId", record.id);
                    }
                  }}
                />
              </div>
            </div>
          </label>
          <Form.Item name="caProviderAccessId" rules={[formRule]}>
            <AccessSelect
              placeholder={t("workflow_node.apply.form.ca_provider_access.placeholder")}
              showSearch
              onFilter={(_, option) => {
                if (option.reserve !== "ca") return false;
                if (fieldCAProvider) return caProvidersMap.get(fieldCAProvider)?.provider === option.provider;

                const provider = accessProvidersMap.get(option.provider);
                return !!provider?.usages?.includes(ACCESS_USAGES.CA);
              }}
            />
          </Form.Item>
        </Form.Item>

        <Form.Item name="keyAlgorithm" label={t("workflow_node.apply.form.key_algorithm.label")} rules={[formRule]}>
          <Select
            options={["RSA2048", "RSA3072", "RSA4096", "RSA8192", "EC256", "EC384"].map((e) => ({
              label: e,
              value: e,
            }))}
            placeholder={t("workflow_node.apply.form.key_algorithm.placeholder")}
          />
        </Form.Item>

        <Form.Item
          name="acmeProfile"
          label={t("workflow_node.apply.form.acme_profile.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.acme_profile.tooltip") }}></span>}
        >
          <AutoComplete
            allowClear
            options={["classic", "tlsserver", "shortlived"].map((value) => ({ value }))}
            placeholder={t("workflow_node.apply.form.acme_profile.placeholder")}
            filterOption={(inputValue, option) => option!.value.toLowerCase().includes(inputValue.toLowerCase())}
          />
        </Form.Item>

        <Divider size="small">
          <Typography.Text className="text-xs font-normal" type="secondary">
            {t("workflow_node.apply.form.advanced_config.label")}
          </Typography.Text>
        </Divider>

        <Form.Item
          name="nameservers"
          label={t("workflow_node.apply.form.nameservers.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.nameservers.tooltip") }}></span>}
        >
          <MultipleSplitValueInput
            modalTitle={t("workflow_node.apply.form.nameservers.multiple_input_modal.title")}
            placeholder={t("workflow_node.apply.form.nameservers.placeholder")}
            placeholderInModal={t("workflow_node.apply.form.nameservers.multiple_input_modal.placeholder")}
            separator={MULTIPLE_INPUT_SEPARATOR}
            splitOptions={{ removeEmpty: true, trimSpace: true }}
          />
        </Form.Item>

        <Form.Item
          name="dnsPropagationWait"
          label={t("workflow_node.apply.form.dns_propagation_wait.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.dns_propagation_wait.tooltip") }}></span>}
        >
          <Input
            type="number"
            allowClear
            min={0}
            max={3600}
            placeholder={t("workflow_node.apply.form.dns_propagation_wait.placeholder")}
            addonAfter={t("workflow_node.apply.form.dns_propagation_wait.unit")}
          />
        </Form.Item>

        <Form.Item
          name="dnsPropagationTimeout"
          label={t("workflow_node.apply.form.dns_propagation_timeout.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.dns_propagation_timeout.tooltip") }}></span>}
        >
          <Input
            type="number"
            allowClear
            min={0}
            max={3600}
            placeholder={t("workflow_node.apply.form.dns_propagation_timeout.placeholder")}
            addonAfter={t("workflow_node.apply.form.dns_propagation_timeout.unit")}
          />
        </Form.Item>

        <Form.Item
          name="dnsTTL"
          label={t("workflow_node.apply.form.dns_ttl.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.dns_ttl.tooltip") }}></span>}
        >
          <Input
            type="number"
            allowClear
            min={0}
            max={86400}
            placeholder={t("workflow_node.apply.form.dns_ttl.placeholder")}
            addonAfter={t("workflow_node.apply.form.dns_ttl.unit")}
          />
        </Form.Item>

        <Form.Item
          name="disableFollowCNAME"
          label={t("workflow_node.apply.form.disable_follow_cname.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.disable_follow_cname.tooltip") }}></span>}
        >
          <Switch />
        </Form.Item>

        <Form.Item
          name="disableARI"
          label={t("workflow_node.apply.form.disable_ari.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.disable_ari.tooltip") }}></span>}
        >
          <Switch />
        </Form.Item>

        <Divider size="small">
          <Typography.Text className="text-xs font-normal" type="secondary">
            {t("workflow_node.apply.form.strategy_config.label")}
          </Typography.Text>
        </Divider>

        <Form.Item
          label={t("workflow_node.apply.form.skip_before_expiry_days.label")}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.skip_before_expiry_days.tooltip") }}></span>}
        >
          <Flex align="center" gap={8} wrap="wrap">
            <div>{t("workflow_node.apply.form.skip_before_expiry_days.prefix")}</div>
            <Form.Item name="skipBeforeExpiryDays" noStyle rules={[formRule]}>
              <InputNumber
                className="w-24"
                min={1}
                max={365}
                placeholder={t("workflow_node.apply.form.skip_before_expiry_days.placeholder")}
                addonAfter={t("workflow_node.apply.form.skip_before_expiry_days.unit")}
              />
            </Form.Item>
            <div>{t("workflow_node.apply.form.skip_before_expiry_days.suffix")}</div>
          </Flex>
        </Form.Item>
      </Form>
    </NodeFormContextProvider>
  );
};

const EmailInput = memo(
  ({ disabled, placeholder, ...props }: { disabled?: boolean; placeholder?: string; value?: string; onChange?: (value: string) => void }) => {
    const { emails, fetchEmails, removeEmail } = useContactEmailsStore();
    useEffect(() => {
      fetchEmails();
    }, []);

    const [value, setValue] = useControllableValue<string>(props, {
      valuePropName: "value",
      defaultValuePropName: "defaultValue",
      trigger: "onChange",
    });

    const [inputValue, setInputValue] = useState<string>();

    const renderOptionLabel = (email: string, removable: boolean = false) => (
      <div className="flex items-center gap-2 overflow-hidden">
        <span className="flex-1 truncate overflow-hidden">{email}</span>
        {removable && (
          <Button
            color="default"
            disabled={disabled}
            icon={<IconCircleMinus size="1.25em" />}
            size="small"
            type="text"
            onClick={(e) => {
              removeEmail(email);
              e.stopPropagation();
            }}
          />
        )}
      </div>
    );

    const options = useMemo(() => {
      const temp = emails.map((email) => ({
        label: renderOptionLabel(email, true),
        value: email,
      }));

      if (!!inputValue && temp.every((option) => option.value !== inputValue)) {
        temp.unshift({
          label: renderOptionLabel(inputValue),
          value: inputValue,
        });
      }

      return temp;
    }, [emails, inputValue]);

    const handleChange = (value: string) => {
      setValue(value);
    };

    const handleSearch = (value: string) => {
      setInputValue(value?.trim());
    };

    return (
      <AutoComplete
        backfill
        defaultValue={value}
        disabled={disabled}
        filterOption
        options={options}
        placeholder={placeholder}
        showSearch
        value={value}
        onChange={handleChange}
        onSearch={handleSearch}
      />
    );
  }
);

const getInitialValues = (): Nullish<z.infer<ReturnType<typeof getSchema>>> => {
  return {
    domains: "",
    provider: "",
    providerAccessId: "",
    ...defaultNodeConfigForApply(),
  };
};

const getSchema = ({ i18n = getI18n() }: { i18n: ReturnType<typeof getI18n> }) => {
  const { t } = i18n;

  return z
    .object({
      domains: z.string().refine((v) => {
        if (!v) return false;
        return String(v)
          .split(MULTIPLE_INPUT_SEPARATOR)
          .every((e) => validDomainName(e, { allowWildcard: true }));
      }, t("common.errmsg.domain_invalid")),
      contactEmail: z.email(t("common.errmsg.email_invalid")),
      challengeType: z.string().nullish(),
      provider: z.string().nonempty(t("workflow_node.apply.form.provider.placeholder")),
      providerAccessId: z.string().nonempty(t("workflow_node.apply.form.provider_access.placeholder")),
      providerConfig: z.any().nullish(),
      caProvider: z.string().nullish(),
      caProviderAccessId: z.string().nullish(),
      caProviderConfig: z.any().nullish(),
      keyAlgorithm: z.string().nonempty(t("workflow_node.apply.form.key_algorithm.placeholder")),
      nameservers: z
        .string()
        .nullish()
        .refine((v) => {
          if (!v) return true;

          return String(v)
            .split(MULTIPLE_INPUT_SEPARATOR)
            .every((e) => validIPv4Address(e) || validIPv6Address(e) || validDomainName(e));
        }, t("common.errmsg.host_invalid")),
      dnsPropagationWait: z.preprocess(
        (v) => (v == null || v === "" ? void 0 : Number(v)),
        z
          .number()
          .int(t("workflow_node.apply.form.dns_propagation_wait.placeholder"))
          .gte(0, t("workflow_node.apply.form.dns_propagation_wait.placeholder"))
          .nullish()
      ),
      dnsPropagationTimeout: z.preprocess(
        (v) => (v == null || v === "" ? void 0 : Number(v)),
        z
          .number()
          .int(t("workflow_node.apply.form.dns_propagation_timeout.placeholder"))
          .gte(1, t("workflow_node.apply.form.dns_propagation_timeout.placeholder"))
          .nullish()
      ),
      dnsTTL: z.preprocess(
        (v) => (v == null || v === "" ? void 0 : Number(v)),
        z.number().int(t("workflow_node.apply.form.dns_ttl.placeholder")).gte(1, t("workflow_node.apply.form.dns_ttl.placeholder")).nullish()
      ),
      acmeProfile: z.string().nullish(),
      disableFollowCNAME: z.boolean().nullish(),
      disableARI: z.boolean().nullish(),
      skipBeforeExpiryDays: z.preprocess(
        (v) => Number(v),
        z
          .number()
          .int(t("workflow_node.apply.form.skip_before_expiry_days.placeholder"))
          .gte(1, t("workflow_node.apply.form.skip_before_expiry_days.placeholder"))
      ),
    })
    .superRefine((values, ctx) => {
      if (values.caProvider) {
        const provider = caProvidersMap.get(values.caProvider);
        if (!provider?.builtin && !values.caProviderAccessId) {
          ctx.addIssue({
            code: "custom",
            message: t("workflow_node.apply.form.ca_provider_access.placeholder"),
            path: ["caProviderAccessId"],
          });
        }
      }
    });
};

const _default = Object.assign(BizApplyNodeConfigForm, {
  getSchema,
});

export default _default;
