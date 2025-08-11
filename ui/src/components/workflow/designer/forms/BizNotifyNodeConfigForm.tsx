import { useEffect, useMemo, useState } from "react";
import { getI18n, useTranslation } from "react-i18next";
import { type FlowNodeEntity, getNodeForm } from "@flowgram.ai/fixed-layout-editor";
import { IconPlus } from "@tabler/icons-react";
import { type AnchorProps, Button, Divider, Flex, Form, type FormInstance, Input, Switch, Typography } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import AccessEditDrawer from "@/components/access/AccessEditDrawer";
import AccessSelect from "@/components/access/AccessSelect";
import NotificationProviderSelect from "@/components/provider/NotificationProviderSelect";
import { ACCESS_USAGES, NOTIFICATION_PROVIDERS, accessProvidersMap, notificationProvidersMap } from "@/domain/provider";
import { type WorkflowNodeConfigForNotify, defaultNodeConfigForNotify } from "@/domain/workflow";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { useAccessesStore } from "@/stores/access";

import { FormNestedFieldsContextProvider, NodeFormContextProvider } from "./_context";
import BizNotifyNodeConfigFormProviderDiscordBot from "./BizNotifyNodeConfigFormProviderDiscordBot";
import BizNotifyNodeConfigFormProviderEmail from "./BizNotifyNodeConfigFormProviderEmail";
import BizNotifyNodeConfigFormProviderMattermost from "./BizNotifyNodeConfigFormProviderMattermost";
import BizNotifyNodeConfigFormProviderSlackBot from "./BizNotifyNodeConfigFormProviderSlackBot";
import BizNotifyNodeConfigFormProviderTelegramBot from "./BizNotifyNodeConfigFormProviderTelegramBot";
import BizNotifyNodeConfigFormProviderWebhook from "./BizNotifyNodeConfigFormProviderWebhook";
import { NodeType } from "../nodes/typings";

export interface BizNotifyNodeConfigFormProps {
  form: FormInstance;
  node: FlowNodeEntity;
}

const BizNotifyNodeConfigForm = ({ node, ...props }: BizNotifyNodeConfigFormProps) => {
  if (node.flowNodeType !== NodeType.BizNotify) {
    console.warn(`[certimate] current workflow node type is not: ${NodeType.BizNotify}`);
  }

  const { i18n, t } = useTranslation();

  const { accesses } = useAccessesStore(useZustandShallowSelector("accesses"));

  const initialValues = useMemo(() => {
    return getNodeForm(node)?.getValueIn("config") as WorkflowNodeConfigForNotify | undefined;
  }, [node]);

  const formSchema = getSchema({ i18n });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    form: props.form,
    name: "workflowNodeBizNotifyNodeConfigForm",
    initialValues: initialValues ?? getInitialValues(),
  });

  const fieldProvider = Form.useWatch<string>("provider", { form: formInst, preserve: true });
  const fieldProviderAccessId = Form.useWatch<string>("providerAccessId", { form: formInst, preserve: true });

  const NestedProviderConfigFields = useMemo(() => {
    /*
      注意：如果追加新的子组件，请保持以 ASCII 排序。
      NOTICE: If you add new child component, please keep ASCII order.
      */
    switch (fieldProvider) {
      case NOTIFICATION_PROVIDERS.DISCORDBOT:
        return BizNotifyNodeConfigFormProviderDiscordBot;
      case NOTIFICATION_PROVIDERS.EMAIL:
        return BizNotifyNodeConfigFormProviderEmail;
      case NOTIFICATION_PROVIDERS.MATTERMOST:
        return BizNotifyNodeConfigFormProviderMattermost;
      case NOTIFICATION_PROVIDERS.SLACKBOT:
        return BizNotifyNodeConfigFormProviderSlackBot;
      case NOTIFICATION_PROVIDERS.TELEGRAMBOT:
        return BizNotifyNodeConfigFormProviderTelegramBot;
      case NOTIFICATION_PROVIDERS.WEBHOOK:
        return BizNotifyNodeConfigFormProviderWebhook;
    }
  }, [fieldProvider]);

  const [showProvider, setShowProvider] = useState(false);
  useEffect(() => {
    // 通常情况下每个授权信息只对应一个消息通知提供商，此时无需显示消息通知提供商字段；
    // 如果对应多个，则显示。
    if (fieldProviderAccessId) {
      const access = accesses.find((e) => e.id === fieldProviderAccessId);
      const providers = Array.from(notificationProvidersMap.values()).filter((e) => e.provider === access?.provider);
      setShowProvider(providers.length > 1);
    } else {
      setShowProvider(false);
    }
  }, [accesses, fieldProviderAccessId]);

  const handleProviderSelect = (value: string) => {
    // 切换消息通知提供商时联动授权信息
    if (initialValues?.provider === value) {
      formInst.setFieldValue("providerAccessId", initialValues?.providerAccessId);
    } else {
      if (notificationProvidersMap.get(fieldProvider)?.provider !== notificationProvidersMap.get(value)?.provider) {
        formInst.setFieldValue("providerAccessId", void 0);
      }
    }
  };

  const handleProviderAccessSelect = (value: string) => {
    // 切换授权信息时联动消息通知提供商
    const access = accesses.find((access) => access.id === value);
    const provider = Array.from(notificationProvidersMap.values()).find((provider) => provider.provider === access?.provider);
    if (fieldProvider !== provider?.type) {
      formInst.setFieldValue("provider", provider?.type);
    }
  };

  return (
    <NodeFormContextProvider value={{ node }}>
      <Form {...formProps} clearOnDestroy={true} form={formInst} layout="vertical" preserve={false} scrollToFirstError>
        <div id="parameters" data-anchor="parameters">
          <Form.Item name="subject" label={t("workflow_node.notify.form.subject.label")} rules={[formRule]}>
            <Input placeholder={t("workflow_node.notify.form.subject.placeholder")} />
          </Form.Item>

          <Form.Item name="message" label={t("workflow_node.notify.form.message.label")} rules={[formRule]}>
            <Input.TextArea autoSize={{ minRows: 3, maxRows: 10 }} placeholder={t("workflow_node.notify.form.message.placeholder")} />
          </Form.Item>

          <Form.Item name="provider" label={t("workflow_node.notify.form.provider.label")} hidden={!showProvider} rules={[formRule]}>
            <NotificationProviderSelect
              disabled={!showProvider}
              placeholder={t("workflow_node.notify.form.provider.placeholder")}
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
                  <span>{t("workflow_node.notify.form.provider_access.label")}</span>
                </div>
                <div className="text-right">
                  <AccessEditDrawer
                    mode="create"
                    trigger={
                      <Button size="small" type="link">
                        {t("workflow_node.notify.form.provider_access.button")}
                        <IconPlus size="1.25em" />
                      </Button>
                    }
                    usage="notification"
                    afterSubmit={(record) => {
                      const provider = accessProvidersMap.get(record.provider);
                      if (provider?.usages?.includes(ACCESS_USAGES.NOTIFICATION)) {
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
                placeholder={t("workflow_node.notify.form.provider_access.placeholder")}
                showSearch
                onChange={handleProviderAccessSelect}
                onFilter={(_, option) => {
                  if (option.reserve !== "notification") return false;

                  const provider = accessProvidersMap.get(option.provider);
                  return !!provider?.usages?.includes(ACCESS_USAGES.NOTIFICATION);
                }}
              />
            </Form.Item>
          </Form.Item>

          <FormNestedFieldsContextProvider value={{ parentNamePath: "providerConfig" }}>
            {NestedProviderConfigFields && <NestedProviderConfigFields />}
          </FormNestedFieldsContextProvider>
        </div>

        <div id="strategy" data-anchor="strategy">
          <Divider size="small">
            <Typography.Text className="text-xs font-normal" type="secondary">
              {t("workflow_node.notify.form_anchor.strategy.title")}
            </Typography.Text>
          </Divider>

          <Form.Item label={t("workflow_node.notify.form.skip_on_all_prev_skipped.label")}>
            <Flex align="center" gap={8} wrap="wrap">
              <div>{t("workflow_node.notify.form.skip_on_all_prev_skipped.prefix")}</div>
              <Form.Item name="skipOnAllPrevSkipped" noStyle rules={[formRule]}>
                <Switch
                  checkedChildren={t("workflow_node.notify.form.skip_on_all_prev_skipped.switch.on")}
                  unCheckedChildren={t("workflow_node.notify.form.skip_on_all_prev_skipped.switch.off")}
                />
              </Form.Item>
              <div>{t("workflow_node.notify.form.skip_on_all_prev_skipped.suffix")}</div>
            </Flex>
          </Form.Item>
        </div>
      </Form>
    </NodeFormContextProvider>
  );
};

const getAnchorItems = ({ i18n = getI18n() }: { i18n: ReturnType<typeof getI18n> }): Required<AnchorProps>["items"] => {
  const { t } = i18n;

  return ["parameters", "strategy"].map((key) => ({
    key: key,
    title: t(`workflow_node.notify.form_anchor.${key}.tab`),
    href: "#" + key,
  }));
};

const getInitialValues = (): Nullish<z.infer<ReturnType<typeof getSchema>>> => {
  return {
    subject: "",
    message: "",
    provider: "",
    providerAccessId: "",
    ...defaultNodeConfigForNotify(),
  };
};

const getSchema = ({ i18n = getI18n() }: { i18n: ReturnType<typeof getI18n> }) => {
  const { t } = i18n;

  return z.object({
    subject: z
      .string()
      .min(1, t("workflow_node.notify.form.subject.placeholder"))
      .max(20480, t("common.errmsg.string_max", { max: 20480 })),
    message: z
      .string()
      .min(1, t("workflow_node.notify.form.message.placeholder"))
      .max(20480, t("common.errmsg.string_max", { max: 20480 })),
    provider: z.string().nonempty(t("workflow_node.notify.form.provider.placeholder")),
    providerAccessId: z.string().nonempty(t("workflow_node.notify.form.provider_access.placeholder")),
    providerConfig: z.any().nullish(),
    skipOnAllPrevSkipped: z.boolean().nullish(),
  });
};

const _default = Object.assign(BizNotifyNodeConfigForm, {
  getAnchorItems,
  getSchema,
});

export default _default;
