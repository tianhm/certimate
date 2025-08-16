import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { App, Card, Col, Row, Spin, Typography } from "antd";
import dayjs from "dayjs";

import {
  WORKFLOW_NODE_TYPES,
  type WorkflowModel,
  type WorkflowNodeConfigForBizDeploy,
  type WorkflowNodeConfigForBizNotify,
  type WorkflowNodeConfigForBranchBlock,
  newNode,
} from "@/domain/workflow";
import { save as saveWorkflow } from "@/repository/workflow";
import { getErrMsg } from "@/utils/error";

const TEMPLATE_KEY_STANDARD = "standard" as const;
const TEMPLATE_KEY_CERTTEST = "certtest" as const;
const TEMPLATE_KEY_EMPTY = "empty" as const;
type TemplateKeys = typeof TEMPLATE_KEY_EMPTY | typeof TEMPLATE_KEY_CERTTEST | typeof TEMPLATE_KEY_STANDARD;

const WorkflowNew = () => {
  const navigate = useNavigate();

  const { i18n, t } = useTranslation();

  const { notification } = App.useApp();

  const templateGridSpans = {
    xs: { flex: "100%" },
    md: { flex: "50%" },
    lg: { flex: "50%" },
    xl: { flex: "33.3333%" },
    xxl: { flex: "33.3333%" },
  };
  const [templateSelectKey, setTemplateSelectKey] = useState<TemplateKeys>();

  const [pending, setPending] = useState(false);

  const handleTemplateClick = async (key: TemplateKeys) => {
    if (pending) return;

    setTemplateSelectKey(key);

    try {
      let workflow = {} as WorkflowModel;
      workflow.name = t("workflow.new.templates.default_name");
      workflow.description = t("workflow.new.templates.default_description", { date: dayjs().format("YYYY-MM-DD HH:mm") });
      workflow.graphDraft = { nodes: [] };
      workflow.hasDraft = true;

      switch (key) {
        case TEMPLATE_KEY_EMPTY:
          {
            const startNode = newNode(WORKFLOW_NODE_TYPES.START, { i18n: i18n });
            const endNode = newNode(WORKFLOW_NODE_TYPES.END, { i18n: i18n });

            workflow.graphDraft!.nodes = [startNode, endNode];
          }
          break;

        case TEMPLATE_KEY_STANDARD:
          {
            const startNode = newNode(WORKFLOW_NODE_TYPES.START, { i18n: i18n });
            const tryCatchNode = newNode(WORKFLOW_NODE_TYPES.TRYCATCH, { i18n: i18n });
            const applyNode = newNode(WORKFLOW_NODE_TYPES.BIZ_APPLY, { i18n: i18n });
            const deployNode = newNode(WORKFLOW_NODE_TYPES.BIZ_DEPLOY, { i18n: i18n });
            const notifyOnSuccessNode = newNode(WORKFLOW_NODE_TYPES.BIZ_NOTIFY, { i18n: i18n });
            const notifyOnFailureNode = newNode(WORKFLOW_NODE_TYPES.BIZ_NOTIFY, { i18n: i18n });
            const endNode = newNode(WORKFLOW_NODE_TYPES.END, { i18n: i18n });

            deployNode.data.config = {
              ...deployNode.data.config,
              certificateOutputNodeId: applyNode.id,
            } as WorkflowNodeConfigForBizDeploy;

            notifyOnSuccessNode.data.config = {
              ...notifyOnSuccessNode.data.config,
              subject: "[Certimate] Workflow Complete",
              message: "Your workflow run has completed successfully.",
              skipOnAllPrevSkipped: true,
            } as WorkflowNodeConfigForBizNotify;

            notifyOnFailureNode.data.config = {
              ...notifyOnFailureNode.data.config,
              subject: "[Certimate] Workflow Failure Alert!",
              message: "Your workflow run has failed. Please check the details.",
            } as WorkflowNodeConfigForBizNotify;

            tryCatchNode.blocks!.at(0)!.blocks ??= [];
            tryCatchNode.blocks!.at(0)!.blocks!.push(applyNode, deployNode, notifyOnSuccessNode);
            tryCatchNode.blocks!.at(1)!.blocks ??= [];
            tryCatchNode.blocks!.at(1)!.blocks!.unshift(notifyOnFailureNode);

            workflow.graphDraft!.nodes = [startNode, tryCatchNode, endNode];
          }
          break;

        case TEMPLATE_KEY_CERTTEST:
          {
            const startNode = newNode(WORKFLOW_NODE_TYPES.START, { i18n: i18n });
            const tryCatchNode = newNode(WORKFLOW_NODE_TYPES.TRYCATCH, { i18n: i18n });
            const monitorNode = newNode(WORKFLOW_NODE_TYPES.BIZ_MONITOR, { i18n: i18n });
            const conditionNode = newNode(WORKFLOW_NODE_TYPES.CONDITION, { i18n: i18n });
            const notifyOnExpireSoonNode = newNode(WORKFLOW_NODE_TYPES.BIZ_NOTIFY, { i18n: i18n });
            const notifyOnExpiredNode = newNode(WORKFLOW_NODE_TYPES.BIZ_NOTIFY, { i18n: i18n });
            const notifyOnFailureNode = newNode(WORKFLOW_NODE_TYPES.BIZ_NOTIFY, { i18n: i18n });
            const endNode = newNode(WORKFLOW_NODE_TYPES.END, { i18n: i18n });

            notifyOnExpireSoonNode.data.config = {
              ...notifyOnExpireSoonNode.data.config,
              subject: "[Certimate] Certificate Expiry Alert!",
              message: "The certificate which you are monitoring will expire soon. Please pay attention to your website.",
            } as WorkflowNodeConfigForBizNotify;

            notifyOnExpiredNode.data.config = {
              ...notifyOnExpiredNode.data.config,
              subject: "[Certimate] Certificate Expiry Alert!",
              message: "The certificate which you are monitoring has already expired. Please pay attention to your website.",
            } as WorkflowNodeConfigForBizNotify;

            notifyOnFailureNode.data.config = {
              ...notifyOnFailureNode.data.config,
              subject: "[Certimate] Workflow Failure Alert!",
              message: "Your workflow run has failed. Please check the details.",
            } as WorkflowNodeConfigForBizNotify;

            tryCatchNode.blocks!.at(0)!.blocks ??= [];
            tryCatchNode.blocks!.at(0)!.blocks!.push(monitorNode, conditionNode);
            tryCatchNode.blocks!.at(1)!.blocks ??= [];
            tryCatchNode.blocks!.at(1)!.blocks!.unshift(notifyOnFailureNode);

            conditionNode.blocks!.at(0)!.data.name = t("workflow_node.condition.default_name.template_certtest_on_expire_soon");
            conditionNode.blocks!.at(0)!.data.config = {
              ...conditionNode.blocks!.at(0)!.data.config,
              expression: {
                left: {
                  left: {
                    selector: {
                      id: monitorNode.id,
                      name: "certificate.validity",
                      type: "boolean",
                    },
                    type: "var",
                  },
                  operator: "eq",
                  right: {
                    type: "const",
                    value: "true",
                    valueType: "boolean",
                  },
                  type: "comparison",
                },
                operator: "and",
                right: {
                  left: {
                    selector: {
                      id: monitorNode.id,
                      name: "certificate.daysLeft",
                      type: "number",
                    },
                    type: "var",
                  },
                  operator: "lte",
                  right: {
                    type: "const",
                    value: "30",
                    valueType: "number",
                  },
                  type: "comparison",
                },
                type: "logical",
              },
            } as WorkflowNodeConfigForBranchBlock;
            conditionNode.blocks!.at(0)!.blocks ??= [];
            conditionNode.blocks!.at(0)!.blocks!.push(notifyOnExpireSoonNode);
            conditionNode.blocks!.at(1)!.data.name = t("workflow_node.condition.default_name.template_certtest_on_expired");
            conditionNode.blocks!.at(1)!.data.config = {
              ...conditionNode.blocks!.at(1)!.data.config,
              expression: {
                left: {
                  selector: {
                    id: monitorNode.id,
                    name: "certificate.validity",
                    type: "boolean",
                  },
                  type: "var",
                },
                operator: "eq",
                right: {
                  type: "const",
                  value: "false",
                  valueType: "boolean",
                },
                type: "comparison",
              },
            } as WorkflowNodeConfigForBranchBlock;
            conditionNode.blocks!.at(1)!.blocks ??= [];
            conditionNode.blocks!.at(1)!.blocks!.push(notifyOnExpiredNode);

            workflow.graphDraft!.nodes = [startNode, tryCatchNode, endNode];
          }
          break;

        default:
          throw "Invalid value of `templateSelectKey`";
      }

      workflow = await saveWorkflow(workflow);
      navigate(`/workflows/${workflow.id}`, { replace: true });
    } catch (err) {
      notification.error({ message: t("common.text.request_error"), description: getErrMsg(err) });

      throw err;
    } finally {
      setPending(false);
      setTemplateSelectKey(void 0);
    }
  };

  return (
    <div className="px-6 py-4">
      <div className="container">
        <h1>{t("workflow.new.title")}</h1>
        <p className="text-base text-gray-500">{t("workflow.new.subtitle")}</p>
      </div>

      <div className="container">
        <Typography.Text type="secondary">
          <div className="mb-4 text-xl">{t("workflow.new.templates.title")}</div>
        </Typography.Text>

        <Row className="justify-stretch" gutter={[16, 16]}>
          <Col {...templateGridSpans}>
            <Card
              className="size-full"
              cover={<img className="min-h-[120px] object-contain" src="/imgs/workflow/tpl-standard.png" />}
              hoverable
              onClick={() => handleTemplateClick(TEMPLATE_KEY_STANDARD)}
            >
              <div className="flex w-full items-center gap-4">
                <Card.Meta
                  className="grow"
                  title={t("workflow.new.templates.template.standard.title")}
                  description={t("workflow.new.templates.template.standard.description")}
                />
                <Spin spinning={templateSelectKey === TEMPLATE_KEY_STANDARD} />
              </div>
            </Card>
          </Col>

          <Col {...templateGridSpans}>
            <Card
              className="size-full"
              cover={<img className="min-h-[120px] object-contain" src="/imgs/workflow/tpl-certtest.png" />}
              hoverable
              onClick={() => handleTemplateClick(TEMPLATE_KEY_CERTTEST)}
            >
              <div className="flex w-full items-center gap-4">
                <Card.Meta
                  className="grow"
                  title={t("workflow.new.templates.template.certtest.title")}
                  description={t("workflow.new.templates.template.certtest.description")}
                />
                <Spin spinning={templateSelectKey === TEMPLATE_KEY_CERTTEST} />
              </div>
            </Card>
          </Col>

          <Col {...templateGridSpans}>
            <Card
              className="size-full"
              cover={<img className="min-h-[120px] object-contain" src="/imgs/workflow/tpl-blank.png" />}
              hoverable
              onClick={() => handleTemplateClick(TEMPLATE_KEY_EMPTY)}
            >
              <div className="flex w-full items-center gap-4">
                <Card.Meta
                  className="grow"
                  title={t("workflow.new.templates.template.empty.title")}
                  description={t("workflow.new.templates.template.empty.description")}
                />
                <Spin spinning={templateSelectKey === TEMPLATE_KEY_EMPTY} />
              </div>
            </Card>
          </Col>
        </Row>
      </div>
    </div>
  );
};

export default WorkflowNew;
