import { getI18n } from "react-i18next";
import { FeedbackLevel, Field } from "@flowgram.ai/fixed-layout-editor";
import { IconCloudUpload } from "@tabler/icons-react";

import { defaultNodeConfigForUpload } from "@/domain/workflow";

import { BaseNode } from "./_shared";
import { newNodeId } from "../_util";
import { NodeKindType, type NodeRegistry, NodeType } from "./typings";
import BizUploadNodeConfigForm from "../forms/BizUploadNodeConfigForm";

export const BizUploadNodeRegistry: NodeRegistry = {
  type: NodeType.BizUpload,
  kindType: NodeKindType.Business,

  meta: {
    helpText: getI18n().t("workflow_node.upload.help"),
    labelText: getI18n().t("workflow_node.upload.label"),

    icon: IconCloudUpload,
    iconColor: "#fff",
    iconBgColor: "#5b65f5",

    clickable: true,
  },

  formMeta: {
    validate: {
      ["config"]: ({ value }) => {
        const res = BizUploadNodeConfigForm.getSchema({}).safeParse(value);
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
          description={<Field<string> name="config.domains">{({ field: { value } }) => <>{value || t("workflow.detail.design.editor.placeholder")}</>}</Field>}
        />
      );
    },
  },

  onAdd: () => {
    const { t } = getI18n();

    return {
      id: newNodeId(),
      type: NodeType.BizUpload,
      data: {
        name: t("workflow_node.upload.default_name"),
        config: defaultNodeConfigForUpload(),
      },
    };
  },
};
