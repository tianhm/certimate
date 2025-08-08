import { useState } from "react";
import { IconX } from "@tabler/icons-react";
import { useControllableValue } from "ahooks";
import { Button, Drawer, Flex } from "antd";

import Show from "@/components/Show";
import { type CertificateModel } from "@/domain/certificate";
import { useTriggerElement } from "@/hooks";
import CertificateDetail from "./CertificateDetail";

export interface CertificateDetailDrawerProps {
  data?: CertificateModel;
  loading?: boolean;
  open?: boolean;
  trigger?: React.ReactNode;
  onOpenChange?: (open: boolean) => void;
}

const CertificateDetailDrawer = ({ data, loading, trigger, ...props }: CertificateDetailDrawerProps) => {
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
        afterOpenChange={setOpen}
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

const useProps = () => {
  const [data, setData] = useState<CertificateDetailDrawerProps["data"]>();
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

const _default = Object.assign(CertificateDetailDrawer, {
  useProps,
});

export default _default;
