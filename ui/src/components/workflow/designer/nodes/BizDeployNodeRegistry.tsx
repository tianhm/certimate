import { getI18n } from "react-i18next";
import { FeedbackLevel, Field } from "@flowgram.ai/fixed-layout-editor";
import { IconPackage } from "@tabler/icons-react";
import { Avatar } from "antd";
import { nanoid } from "nanoid";

import { deploymentProvidersMap } from "@/domain/provider";

import { BaseNode } from "./_shared";
import { type NodeRegistry, NodeType } from "./typings";

export const BizDeployNodeRegistry: NodeRegistry = {
  type: NodeType.BizDeploy,

  meta: {
    helpText: getI18n().t("workflow_node.deploy.help"),
    labelText: getI18n().t("workflow_node.deploy.label"),

    icon: IconPackage,
    iconColor: "#fff",
    iconBgColor: "#5b65f5",

    clickable: true,
  },

  formMeta: {
    validate: {
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
              <Field<string> name="config.provider">
                {({ field: { value } }) => (
                  <>
                    {value ? (
                      <>
                        <div className="flex-1 truncate">{t(deploymentProvidersMap.get(value)?.name ?? "")}</div>
                        <Avatar shape="square" src={deploymentProvidersMap.get(value)?.icon} size={20} />
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
      type: NodeType.BizMonitor,
      data: {
        name: t("workflow_node.deploy.default_name"),
      },
    };
  },
};
