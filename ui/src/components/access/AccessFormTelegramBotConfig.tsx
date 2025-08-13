import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AccessConfigForTelegramBot } from "@/domain/access";

type AccessFormTelegramBotConfigFieldValues = Nullish<AccessConfigForTelegramBot>;

export interface AccessFormTelegramBotConfigProps {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessFormTelegramBotConfigFieldValues;
  onValuesChange?: (values: AccessFormTelegramBotConfigFieldValues) => void;
}

const initFormModel = (): AccessFormTelegramBotConfigFieldValues => {
  return {
    botToken: "",
  };
};

const AccessFormTelegramBotConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: AccessFormTelegramBotConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    botToken: z
      .string()
      .min(1, t("access.form.telegrambot_token.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 })),
    chatId: z
      .preprocess(
        (v) => (v == null || v === "" ? void 0 : Number(v)),
        z.number().refine((v) => {
          return !Number.isNaN(+v!) && +v! !== 0;
        }, t("access.form.telegrambot_chat_id.placeholder"))
      )
      .nullish(),
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
        name="botToken"
        label={t("access.form.telegrambot_token.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.telegrambot_token.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.telegrambot_token.placeholder")} />
      </Form.Item>

      <Form.Item
        name="chatId"
        label={t("access.form.telegrambot_chat_id.label")}
        extra={t("access.form.telegrambot_chat_id.help")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.telegrambot_chat_id.tooltip") }}></span>}
      >
        <Input allowClear placeholder={t("access.form.telegrambot_chat_id.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessFormTelegramBotConfig;
