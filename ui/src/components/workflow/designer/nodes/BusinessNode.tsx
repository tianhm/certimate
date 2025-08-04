import { getI18n } from "react-i18next";
import { IconCloudUpload, IconContract, IconDeviceDesktopSearch, IconPackage, IconSend } from "@tabler/icons-react";
import { Avatar, Typography } from "antd";
import { nanoid } from "nanoid";

import { deploymentProvidersMap, notificationProvidersMap } from "@/domain/provider";

import { BaseNode } from "./_shared";
import { type NodeRegistry, NodeType } from "./typings";

export const BizApplyNodeRegistry: NodeRegistry = {
  type: NodeType.BizApply,

  meta: {
    helpText: getI18n().t("workflow_node.apply.help"),

    icon: IconContract,
    iconColor: "#fff",
    iconBgColor: "#5b65f5",
  },

  formMeta: {
    render: ({ form }) => {
      const { t } = getI18n();

      const fieldDomains = form.getValueIn<string>("config.domains");

      return <BaseNode>{fieldDomains ? fieldDomains.split(";").join("; ") : t("workflow.detail.design.nodes.placeholder")}</BaseNode>;
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

    icon: IconCloudUpload,
    iconColor: "#fff",
    iconBgColor: "#5b65f5",
  },

  formMeta: {
    render: ({ form }) => {
      const { t } = getI18n();

      const fieldDomains = form.getValueIn<string>("config.domains");

      return <BaseNode>{fieldDomains ? fieldDomains.split(";").join("; ") : t("workflow.detail.design.nodes.placeholder")}</BaseNode>;
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

    icon: IconDeviceDesktopSearch,
    iconColor: "#fff",
    iconBgColor: "#5b65f5",
  },

  formMeta: {
    render: ({ form }) => {
      const { t } = getI18n();

      const fieldDomain = form.getValueIn<string>("config.domain");
      const fieldHost = form.getValueIn<string>("config.host");

      return <BaseNode>{fieldDomain || fieldHost ? fieldDomain || fieldHost : t("workflow.detail.design.nodes.placeholder")}</BaseNode>;
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

    icon: IconPackage,
    iconColor: "#fff",
    iconBgColor: "#5b65f5",
  },

  formMeta: {
    render: ({ form }) => {
      const { t } = getI18n();

      const fieldProvider = form.getValueIn<string>("config.provider");

      return (
        <BaseNode>
          <div className="flex items-center justify-between gap-1">
            {fieldProvider ? (
              <>
                <div className="flex-1 truncate">{t(deploymentProvidersMap.get(fieldProvider)?.name ?? "")}</div>
                <Avatar shape="square" src={deploymentProvidersMap.get(fieldProvider)?.icon} size={20} />
              </>
            ) : (
              t("workflow.detail.design.nodes.placeholder")
            )}
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

    icon: IconSend,
    iconColor: "#fff",
    iconBgColor: "#0693d4",
  },

  formMeta: {
    render: ({ form }) => {
      const { t } = getI18n();

      const fieldProvider = form.getValueIn<string>("config.provider");

      return (
        <BaseNode>
          <div className="flex items-center justify-between gap-1">
            {fieldProvider ? (
              <>
                <div className="flex-1 truncate">{t(notificationProvidersMap.get(fieldProvider)?.name ?? "")}</div>
                <Avatar shape="square" src={notificationProvidersMap.get(fieldProvider)?.icon} size="small" />
              </>
            ) : (
              t("workflow.detail.design.nodes.placeholder")
            )}
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
