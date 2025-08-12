import { getI18n } from "react-i18next";
import { FeedbackLevel, Field } from "@flowgram.ai/fixed-layout-editor";
import { IconSend } from "@tabler/icons-react";
import { Avatar } from "antd";
import { nanoid } from "nanoid";

import { notificationProvidersMap } from "@/domain/provider";

import { BaseNode } from "./_shared";
import { NodeKindType, type NodeRegistry, NodeType } from "./typings";
import BizNotifyNodeConfigForm from "../forms/BizNotifyNodeConfigForm";

export const BizNotifyNodeRegistry: NodeRegistry = {
  type: NodeType.BizNotify,
  kindType: NodeKindType.Business,

  meta: {
    helpText: getI18n().t("workflow_node.notify.help"),
    labelText: getI18n().t("workflow_node.notify.label"),

    icon: IconSend,
    iconColor: "#fff",
    iconBgColor: "#0693d4",

    clickable: true,
  },

  formMeta: {
    validate: {
      ["config"]: ({ value }) => {
        const res = BizNotifyNodeConfigForm.getSchema({}).safeParse(value);
        if (!res.success) {
          return {
            message: res.error.message,
            level: FeedbackLevel.Error,
          };
        }
      },
    },

    render: () => {
      const { t } = getI18n();

      return (
        <BaseNode
          description={
            <div className="flex items-center justify-between gap-1">
              <Field<string> name="config.provider">
                {({ field: { value } }) => (
                  <>
                    {value ? (
                      <>
                        <div className="flex-1 truncate">{t(notificationProvidersMap.get(value)?.name ?? "")}</div>
                        <Avatar shape="square" src={notificationProvidersMap.get(value)?.icon} size={20} />
                      </>
                    ) : (
                      t("workflow.detail.design.editor.placeholder")
                    )}
                  </>
                )}
              </Field>
            </div>
          }
        />
      );
    },
  },

  onAdd: () => {
    const { t } = getI18n();

    return {
      id: nanoid(),
      type: NodeType.BizNotify,
      data: {
        name: t("workflow_node.notify.default_name"),
      },
    };
  },
};
