import { getI18n } from "react-i18next";
import { IconRocket } from "@tabler/icons-react";
import { Typography } from "antd";

import { WORKFLOW_TRIGGERS } from "@/domain/workflow";

import { BaseNode } from "./_shared";
import { type NodeRegistry, NodeType } from "./typings";

export const StartNodeRegistry: NodeRegistry = {
  type: NodeType.Start,

  meta: {
    helpText: getI18n().t("workflow_node.start.help"),

    icon: IconRocket,
    iconColor: "#fff",
    iconBgColor: "#ed6d0c",

    isStart: true,

    expandable: false,
    selectable: false,

    addDisable: true,
    copyDisable: true,
    deleteDisable: true,
  },

  formMeta: {
    render: ({ form }) => {
      const { t } = getI18n();

      const fieldTrigger = form.getValueIn<string>("config.trigger");
      const fieldTriggerCron = form.getValueIn<string>("config.triggerCron");

      return (
        <BaseNode>
          <div className="flex items-center justify-between gap-1">
            {fieldTrigger ? (
              <>
                <div>
                  {fieldTrigger === WORKFLOW_TRIGGERS.SCHEDULED
                    ? t("workflow.props.trigger.scheduled")
                    : fieldTrigger === WORKFLOW_TRIGGERS.MANUAL
                      ? t("workflow.props.trigger.manual")
                      : "\u00A0"}
                </div>
                <div>{fieldTrigger === WORKFLOW_TRIGGERS.SCHEDULED ? fieldTriggerCron || "\u00A0" : ""}</div>
              </>
            ) : (
              t("workflow.detail.design.nodes.placeholder")
            )}
          </div>
        </BaseNode>
      );
    },
  },

  canAdd: () => {
    return false;
  },

  canDelete: () => {
    return false;
  },
};
