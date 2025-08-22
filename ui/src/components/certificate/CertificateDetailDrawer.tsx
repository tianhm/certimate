import { startTransition, useCallback, useState } from "react";
import { IconX } from "@tabler/icons-react";
import { useControllableValue } from "ahooks";
import { Button, Drawer, Flex } from "antd";

import Show from "@/components/Show";
import { type CertificateModel } from "@/domain/certificate";
import { useTriggerElement } from "@/hooks";
import CertificateDetail from "./CertificateDetail";

export interface CertificateDetailDrawerProps {
  afterClose?: () => void;
  data?: CertificateModel;
  loading?: boolean;
  open?: boolean;
  trigger?: React.ReactNode;
  onOpenChange?: (open: boolean) => void;
}

const CertificateDetailDrawer = ({ afterClose, data, loading, trigger, ...props }: CertificateDetailDrawerProps) => {
  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const triggerEl = useTriggerElement(trigger, { onClick: () => setOpen(true) });

  return (
    <>
      {triggerEl}

      <Drawer
        afterOpenChange={(open) => !open && afterClose?.()}
        autoFocus
        closeIcon={false}
        destroyOnHidden
        open={open}
        loading={loading}
        placement="right"
        size="large"
        title={
          <Flex align="center" justify="space-between" gap="small">
            <div className="flex-1 truncate">{data ? `Certificate #${data.id}` : "Certificate"}</div>
            <Button
              className="ant-drawer-close"
              style={{ marginInline: 0 }}
              icon={<IconX size="1.25em" />}
              size="small"
              type="text"
              onClick={() => setOpen(false)}
            />
          </Flex>
        }
        onClose={() => setOpen(false)}
      >
        <Show when={!!data}>
          <CertificateDetail data={data!} />
        </Show>
      </Drawer>
    </>
  );
};

const useDrawer = () => {
  type DataType = CertificateDetailDrawerProps["data"];
  const [data, setData] = useState<DataType>();
  const [open, setOpen] = useState<boolean>(false);

  const onOpenChange = useCallback((open: boolean) => {
    setOpen(open);
  }, []);

  return {
    drawerProps: {
      afterClose: () => {
        startTransition(() => {
          if (!open) {
            setData(void 0);
          }
        });
      },
      data,
      open,
      onOpenChange,
    },

    open: (data: NonNullable<DataType>) => {
      setData(data);
      setOpen(true);
    },
    close: () => {
      setOpen(false);
    },
  };
};

const _default = Object.assign(CertificateDetailDrawer, {
  useDrawer,
});

export default _default;
