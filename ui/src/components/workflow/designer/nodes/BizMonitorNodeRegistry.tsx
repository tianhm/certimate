import { getI18n } from "react-i18next";
import { FeedbackLevel, Field } from "@flowgram.ai/fixed-layout-editor";
import { IconDeviceDesktopSearch } from "@tabler/icons-react";
import { nanoid } from "nanoid";

import { BaseNode } from "./_shared";
import { type NodeRegistry, NodeType } from "./typings";

export const BizMonitorNodeRegistry: NodeRegistry = {
  type: NodeType.BizMonitor,

  meta: {
    helpText: getI18n().t("workflow_node.monitor.help"),
    labelText: getI18n().t("workflow_node.monitor.label"),

    icon: IconDeviceDesktopSearch,
    iconColor: "#fff",
    iconBgColor: "#5b65f5",

    clickable: true,
  },

  formMeta: {
    validate: {
      ["config.host"]: ({ value }) => {
        if (!value) {
          return {
            message: "required",
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
            <Field name="config.domain">
              {({ field: { value: fieldDomain } }) => (
                <Field name="config.host">
                  {({ field: { value: fieldHost } }) => (
                    <>{fieldDomain || fieldHost ? fieldDomain || fieldHost : t("workflow.detail.design.editor.placeholder")}</>
                  )}
                </Field>
              )}
            </Field>
          }
        />
      );
    },
  },

  onAdd: () => {
    const { t } = getI18n();

    return {
      id: nanoid(),
      type: NodeType.BizMonitor,
      data: {
        name: t("workflow_node.monitor.default_name"),
      },
    };
  },
};
