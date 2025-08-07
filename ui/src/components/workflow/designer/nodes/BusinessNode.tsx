import { getI18n } from "react-i18next";
import { FeedbackLevel, Field } from "@flowgram.ai/fixed-layout-editor";
import { IconCloudUpload, IconContract, IconDeviceDesktopSearch, IconPackage, IconSend } from "@tabler/icons-react";
import { Avatar } from "antd";
import { nanoid } from "nanoid";

import { deploymentProvidersMap, notificationProvidersMap } from "@/domain/provider";

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
        <BaseNode>
          <Field<string> name="config.domains">{({ field: { value } }) => <>{value || t("workflow.detail.design.editor.placeholder")}</>}</Field>
        </BaseNode>
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

export const BizUploadNodeRegistry: NodeRegistry = {
  type: NodeType.BizUpload,

  meta: {
    helpText: getI18n().t("workflow_node.upload.help"),
    labelText: getI18n().t("workflow_node.upload.label"),

    icon: IconCloudUpload,
    iconColor: "#fff",
    iconBgColor: "#5b65f5",
  },

  formMeta: {
    validate: {
      ["config.certificate"]: ({ value }) => {
        if (!value) {
          return {
            message: "required",
            level: FeedbackLevel.Error,
          };
        }
      },
      ["config.privateKey"]: ({ value }) => {
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
        <BaseNode>
          <Field<string> name="config.domains">{({ field: { value } }) => <>{value || t("workflow.detail.design.editor.placeholder")}</>}</Field>
        </BaseNode>
      );
    },
  },

  onAdd: () => {
    const { t } = getI18n();

    return {
      id: nanoid(),
      type: NodeType.BizUpload,
      data: {
        name: t("workflow_node.upload.default_name"),
      },
    };
  },
};

export const BizMonitorNodeRegistry: NodeRegistry = {
  type: NodeType.BizMonitor,

  meta: {
    helpText: getI18n().t("workflow_node.monitor.help"),
    labelText: getI18n().t("workflow_node.monitor.label"),

    icon: IconDeviceDesktopSearch,
    iconColor: "#fff",
    iconBgColor: "#5b65f5",
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

    render: ({ form }) => {
      const { t } = getI18n();

      const fieldDomain = form.getValueIn<string>("config.domain");
      const fieldHost = form.getValueIn<string>("config.host");

      return <BaseNode>{fieldDomain || fieldHost ? fieldDomain || fieldHost : t("workflow.detail.design.editor.placeholder")}</BaseNode>;
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

export const BizDeployNodeRegistry: NodeRegistry = {
  type: NodeType.BizDeploy,

  meta: {
    helpText: getI18n().t("workflow_node.deploy.help"),
    labelText: getI18n().t("workflow_node.deploy.label"),

    icon: IconPackage,
    iconColor: "#fff",
    iconBgColor: "#5b65f5",
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
        <BaseNode>
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
        </BaseNode>
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

export const BizNotifyNodeRegistry: NodeRegistry = {
  type: NodeType.BizNotify,

  meta: {
    helpText: getI18n().t("workflow_node.notify.help"),
    labelText: getI18n().t("workflow_node.notify.label"),

    icon: IconSend,
    iconColor: "#fff",
    iconBgColor: "#0693d4",
  },

  formMeta: {
    validate: {
      ["config.subject"]: ({ value }) => {
        if (!value) {
          return {
            message: "required",
            level: FeedbackLevel.Error,
          };
        }
      },
      ["config.message"]: ({ value }) => {
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
        <BaseNode>
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
        </BaseNode>
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
