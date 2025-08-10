import { type AdderProps as FlowgramAdderProps, useClientContext } from "@flowgram.ai/fixed-layout-editor";

import { IconPlus } from "@tabler/icons-react";
import { Button, Dropdown } from "antd";

import { getFlowNodeRegistries } from "../nodes";

export interface AdderProps extends FlowgramAdderProps {}

const Adder = ({ from, hoverActivated }: AdderProps) => {
  const ctx = useClientContext();
  const { operation, playground } = ctx;

  const menuItems = getFlowNodeRegistries()
    .filter((registry) => {
      if (registry.meta?.addDisable != null) {
        return !registry.meta.addDisable;
      }
      return true;
    })
    .filter((registry) => {
      if (registry.canAdd != null) {
        return registry.canAdd(ctx, from);
      }
      return true;
    })
    .map((registry) => {
      const Icon = registry.meta?.icon;

      return {
        key: registry.type,
        label: registry.meta?.labelText ?? registry.type,
        icon: <span className="anticon scale-125">{Icon && <Icon size="1em" />}</span>,
        onClick: () => {
          const block = operation.addFromNode(from, registry.onAdd!(ctx, from));

          setTimeout(() => {
            playground.scrollToView({
              bounds: block.bounds,
              scrollToCenter: true,
            });
          }, 1);
        },
      };
    });

  return playground.config.readonlyOrDisabled ? null : (
    <div className="relative">
      <Dropdown menu={{ items: menuItems }} placement="bottomRight" trigger={["click"]}>
        {hoverActivated ? (
          <Button icon={<IconPlus size="1em" stroke={3} />} shape="circle" size="small" type="primary" />
        ) : (
          <div className="size-2 rounded-full bg-primary opacity-75"></div>
        )}
      </Dropdown>
    </div>
  );
};

export default Adder;
