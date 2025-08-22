import { getI18n } from "react-i18next";
import { FeedbackLevel, Field } from "@flowgram.ai/fixed-layout-editor";
import { IconCloudUpload } from "@tabler/icons-react";

import { newNode } from "@/domain/workflow";

import { BaseNode } from "./_shared";
import { NodeKindType, type NodeRegistry, NodeType } from "./typings";
import BizUploadNodeConfigForm from "../forms/BizUploadNodeConfigForm";

export const BizUploadNodeRegistry: NodeRegistry = {
  type: NodeType.BizUpload,

  kind: NodeKindType.Business,

  meta: {
    helpText: getI18n().t("workflow_node.upload.help"),
    labelText: getI18n().t("workflow_node.upload.label"),

    icon: IconCloudUpload,
    iconColor: "#fff",
    iconBgColor: "#5b65f5",

    clickable: true,
    expandable: false,
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
    return newNode(NodeType.BizUpload, { i18n: getI18n() });
  },
};
