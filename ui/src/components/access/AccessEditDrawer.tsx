import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { IconX } from "@tabler/icons-react";
import { useControllableValue } from "ahooks";
import { App, Button, Drawer, Flex } from "antd";

import { type AccessModel } from "@/domain/access";
import { useTriggerElement, useZustandShallowSelector } from "@/hooks";
import { useAccessesStore } from "@/stores/access";
import { getErrMsg } from "@/utils/error";

import AccessForm, { type AccessFormInstance, type AccessFormProps } from "./AccessForm";

export interface AccessEditDrawerProps {
  data?: AccessFormProps["initialValues"];
  loading?: boolean;
  mode: AccessFormProps["mode"];
  open?: boolean;
  trigger?: React.ReactNode;
  usage?: AccessFormProps["usage"];
  onOpenChange?: (open: boolean) => void;
  afterSubmit?: (record: AccessModel) => void;
}

const AccessEditDrawer = ({ mode, data, loading, trigger, usage, afterSubmit, ...props }: AccessEditDrawerProps) => {
  const { t } = useTranslation();

  const { notification } = App.useApp();

  const { createAccess, updateAccess } = useAccessesStore(useZustandShallowSelector(["createAccess", "updateAccess"]));

  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const triggerEl = useTriggerElement(trigger, { onClick: () => setOpen(true) });

  const formRef = useRef<AccessFormInstance>(null);
  const [formPending, setFormPending] = useState(false);

  const [footerShow, setFooterShow] = useState(!!data?.provider);
  useEffect(() => {
    setFooterShow(!!data?.provider);
  }, [data?.provider]);

  const handleOkClick = async () => {
    setFormPending(true);
    try {
      await formRef.current!.validateFields();
    } catch (err) {
      setFormPending(false);
      throw err;
    }

    try {
      let values: AccessModel = formRef.current!.getFieldsValue();

      if (mode === "create") {
        if (data?.id) {
          throw "Invalid props: `data`";
        }

        values = await createAccess(values);
      } else if (mode === "edit") {
        if (!data?.id) {
          throw "Invalid props: `data`";
        }

        values = await updateAccess({ ...data, ...values });
      } else {
        throw "Invalid props: `action`";
      }

      afterSubmit?.(values);
      setOpen(false);
    } catch (err) {
      notification.error({ message: t("common.text.request_error"), description: getErrMsg(err) });

      throw err;
    } finally {
      setFormPending(false);
    }
  };

  const handleCancelClick = () => {
    if (formPending) return;

    setOpen(false);
  };

  const handleFormValuesChange: AccessFormProps["onValuesChange"] = (values) => {
    setFooterShow(!!values.provider);
  };

  return (
    <>
      {triggerEl}

      <Drawer
        afterOpenChange={setOpen}
        closeIcon={false}
        destroyOnHidden
        footer={
          footerShow ? (
            <Flex className="px-2" justify="end" gap="small">
              <Button onClick={handleCancelClick}>{t("common.button.cancel")}</Button>
              <Button loading={formPending} type="primary" onClick={handleOkClick}>
                {mode === "edit" ? t("common.button.save") : t("common.button.submit")}
              </Button>
            </Flex>
          ) : (
            false
          )
        }
        loading={loading}
        maskClosable={!formPending}
        open={open}
        size="large"
        title={
          <Flex align="center" justify="space-between" gap="small">
            <div className="flex-1 truncate">
              {mode === "edit" && !!data ? t("access.action.edit.modal.title") + ` #${data.id}` : t(`access.action.${mode}.modal.title`)}
            </div>
            <Button
              className="ant-drawer-close"
              style={{ marginInline: 0 }}
              icon={<IconX size="1.25em" />}
              size="small"
              type="text"
              onClick={handleCancelClick}
            />
          </Flex>
        }
        onClose={handleCancelClick}
      >
        <AccessForm ref={formRef} disabled={formPending} initialValues={data} mode={mode} usage={usage} onValuesChange={handleFormValuesChange} />
      </Drawer>
    </>
  );
};

const useProps = () => {
  const [data, setData] = useState<AccessEditDrawerProps["data"]>();
  const [open, setOpen] = useState<boolean>(false);

  const onOpenChange = (open: boolean) => {
    setOpen(open);

    if (!open) {
      setData(undefined);
    }
  };

  return {
    data,
    open,
    setData,
    setOpen,
    onOpenChange,
  };
};

const _default = Object.assign(AccessEditDrawer, {
  useProps,
});

export default _default;
