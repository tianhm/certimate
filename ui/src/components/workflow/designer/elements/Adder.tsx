import { useTranslation } from "react-i18next";
import { type AdderProps as FlowgramAdderProps, useClientContext } from "@flowgram.ai/fixed-layout-editor";

import { IconPlus } from "@tabler/icons-react";
import { Button, Dropdown, type MenuProps } from "antd";

import { getAllNodeRegistries } from "../nodes";

export interface AdderProps extends FlowgramAdderProps {}

const Adder = ({ from, hoverActivated }: AdderProps) => {
  const { t } = useTranslation();

  const ctx = useClientContext();
  const { operation, playground } = ctx;

  const menuItems = getAllNodeRegistries()
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
    .reduce(
      (acc, registry) => {
        let group = acc.find((item) => item!.key === registry.kindType);
        if (!group) {
          group = {
            key: registry.kindType,
            type: "group",
            label: registry.kindType ? t(`workflow_node.kind.${registry.kindType}`) : null,
            children: [],
          };
          acc.push(group);
        }

        if (group.type === "group") {
          const NodeIcon = registry.meta?.icon;
          group.children!.push({
            key: registry.type,
            label: registry.meta?.labelText ?? registry.type,
            icon: <span className="anticon scale-125">{NodeIcon && <NodeIcon size="1em" />}</span>,
            onClick: () => {
              const block = operation.addFromNode(from, registry.onAdd!(ctx, from));

              setTimeout(() => {
                playground.scrollToView({
                  bounds: block.bounds,
                  scrollToCenter: true,
                });
              }, 1);
            },
          });
        }

        return acc;
      },
      [] as Required<MenuProps>["items"]
    );

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
