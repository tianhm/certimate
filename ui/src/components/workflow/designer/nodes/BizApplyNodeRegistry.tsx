import { getI18n } from "react-i18next";
import { FeedbackLevel, Field } from "@flowgram.ai/fixed-layout-editor";
import { IconContract } from "@tabler/icons-react";
import { Avatar } from "antd";
import { nanoid } from "nanoid";

import { acmeDns01ProvidersMap } from "@/domain/provider";

import { BaseNode } from "./_shared";
import { type NodeRegistry, NodeType } from "./typings";

export const BizApplyNodeRegistry: NodeRegistry = {
  type: NodeType.BizApply,

  meta: {
    helpText: getI18n().t("workflow_node.apply.help"),
    labelText: getI18n().t("workflow_node.apply.label"),

    icon: IconContract,
    iconColor: "#fff",
    iconBgColor: "#5b65f5",

    clickable: true,
  },

  formMeta: {
    validate: {
      ["config.domains"]: ({ value }) => {
        if (!value) {
          return {
            message: "required",
            level: FeedbackLevel.Error,
          };
        }
      },
      ["config.contactEmail"]: ({ value }) => {
        if (!value) {
          return {
            message: "required",
            level: FeedbackLevel.Error,
          };
        }
      },
      ["config.provider"]: ({ value }) => {
        if (!value) {
          return {
            message: "required",
            level: FeedbackLevel.Error,
          };
        }
      },
      ["config.providerAccessId"]: ({ value }) => {
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
            <div className="flex items-center justify-between gap-1">
              <Field<string> name="config.domains">
                {({ field: { value } }) => {
                  return <div className="flex-1 truncate">{value || t("workflow.detail.design.editor.placeholder")}</div>;
                }}
              </Field>
              <Field<string> name="config.provider">
                {({ field: { value } }) => (value ? <Avatar shape="square" src={acmeDns01ProvidersMap.get(value)?.icon} size={20} /> : <></>)}
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
      type: NodeType.BizApply,
      data: {
        name: t("workflow_node.apply.default_name"),
      },
    };
  },
};
