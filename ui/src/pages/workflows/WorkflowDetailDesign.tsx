import { useEffect, useState } from "react";
import { getI18n, useTranslation } from "react-i18next";
import { type FlowDocumentJSON } from "@flowgram.ai/document";
import { IconArrowBackUp, IconDots } from "@tabler/icons-react";
import { useDeepCompareEffect } from "ahooks";
import { Alert, App, Button, Card, Dropdown, Space, theme } from "antd";
import { nanoid } from "nanoid";
import { isEqual } from "radash";

import Show from "@/components/Show";
import WorkflowDesigner from "@/components/workflow/designer/Editor";
import WorkflowDesignerToolbar from "@/components/workflow/designer/Toolbar";
import { type WorkflowNode, WorkflowNodeType, isAllNodesValidated } from "@/domain/workflow";
import { WORKFLOW_RUN_STATUSES } from "@/domain/workflowRun";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";
import { getErrMsg } from "@/utils/error";

const WorkflowDetailDesign = () => {
  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();
  const { message, modal, notification } = App.useApp();

  const { workflow, ...workflowState } = useWorkflowStore(useZustandShallowSelector(["workflow", "init", "publish", "rollback"]));

  const [isPendingOrRunning, setIsPendingOrRunning] = useState(false);
  const [allowRollback, setAllowRollback] = useState(false);
  const [allowPublish, setAllowPublish] = useState(false);

  useEffect(() => {
    const pending = workflow.lastRunStatus === WORKFLOW_RUN_STATUSES.PENDING || workflow.lastRunStatus === WORKFLOW_RUN_STATUSES.RUNNING;
    setIsPendingOrRunning(pending);
  }, [workflow]);

  useEffect(() => {
    const hasContent = !!workflow.content;
    const hasChanges = workflow.hasDraft! || !isEqual(workflow.draft, workflow.content);
    setAllowRollback(!isPendingOrRunning && hasContent && hasChanges);
    setAllowPublish(!isPendingOrRunning && hasChanges);
  }, [workflow.content, workflow.draft, workflow.hasDraft, isPendingOrRunning]);

  const [editorData, setEditorData] = useState<FlowDocumentJSON>();
  useDeepCompareEffect(() => {
    setEditorData({ nodes: compactWorkflowDraft(workflow.draft) });
  }, [workflow.draft]);

  const handleRollbackClick = () => {
    modal.confirm({
      title: t("workflow.action.rollback.modal.title"),
      content: t("workflow.action.rollback.modal.content"),
      onOk: async () => {
        try {
          await workflowState.rollback();

          message.success(t("common.text.operation_succeeded"));
        } catch (err) {
          console.error(err);
          notification.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
        }
      },
    });
  };

  const handlePublishClick = () => {
    if (!isAllNodesValidated(workflow.draft!)) {
      message.warning(t("workflow.action.publish.errmsg.uncompleted"));
      return;
    }

    modal.confirm({
      title: t("workflow.action.publish.modal.title"),
      content: t("workflow.action.publish.modal.content"),
      onOk: async () => {
        try {
          await workflowState.publish();

          message.success(t("common.text.operation_succeeded"));
        } catch (err) {
          console.error(err);
          notification.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
        }
      },
    });
  };

  return (
    <div className="size-full">
      <Card
        className="size-full overflow-hidden"
        styles={{
          body: {
            position: "relative",
            height: "100%",
            padding: 0,
          },
        }}
      >
        <WorkflowDesigner initialData={editorData}>
          <div className="absolute top-8 z-10 w-full px-4">
            <div className="container">
              <div className="flex items-center justify-end gap-4">
                <div className="flex flex-1 items-center justify-end gap-4 overflow-hidden">
                  <div className="flex-1 overflow-hidden">
                    <Show when={workflow.hasDraft!}>
                      <Alert message={<div className="truncate">{t("workflow.detail.design.unpublished_draft.alert")}</div>} showIcon type="warning" />
                    </Show>
                  </div>
                  <Space.Compact>
                    <Button disabled={!allowPublish} ghost type="primary" onClick={handlePublishClick}>
                      {t("workflow.action.publish.button")}
                    </Button>
                    <Dropdown
                      menu={{
                        items: [
                          {
                            key: "rollback",
                            disabled: !allowRollback,
                            label: t("workflow.action.rollback.button"),
                            icon: <IconArrowBackUp size="1.25em" />,
                            onClick: handleRollbackClick,
                          },
                        ],
                      }}
                      trigger={["click"]}
                    >
                      <Button icon={<IconDots size="1.25em" />} />
                    </Dropdown>
                  </Space.Compact>
                </div>
              </div>
            </div>
          </div>
          <div className="absolute bottom-8 z-10 w-full px-4">
            <div className="container">
              <div className="flex justify-end">
                <WorkflowDesignerToolbar
                  style={{
                    backgroundColor: themeToken.colorBgContainer,
                    borderRadius: themeToken.borderRadius,
                  }}
                />
              </div>
            </div>
          </div>
        </WorkflowDesigner>
      </Card>
    </div>
  );
};

const compactWorkflowDraft = (root: WorkflowNode | undefined) => {
  const { t } = getI18n();

  // TODO: 仅为兼容适配 v0.3.x 数据，正式上线后待删除
  const res: FlowDocumentJSON["nodes"] = [];

  if (!root) {
    res.push({
      id: nanoid(),
      type: "start",
      data: {
        name: "Start",
      },
    });
  } else {
    const convert = (node: WorkflowNode | undefined) => {
      const temp: typeof res = [];

      let current: typeof node = node;
      while (current) {
        switch (current.type) {
          case WorkflowNodeType.Start:
            temp.push({
              id: current.id,
              type: "start",
              data: {
                name: current.name,
                config: current.config,
              },
            });
            break;

          case WorkflowNodeType.Apply:
            temp.push({
              id: current.id,
              type: "bizApply",
              data: {
                name: current.name,
                config: current.config,
              },
            });
            break;

          case WorkflowNodeType.Upload:
            temp.push({
              id: current.id,
              type: "bizUpload",
              data: {
                name: current.name,
                config: current.config,
              },
            });
            break;

          case WorkflowNodeType.Monitor:
            temp.push({
              id: current.id,
              type: "bizMonitor",
              data: {
                name: current.name,
                config: current.config,
              },
            });
            break;

          case WorkflowNodeType.Deploy:
            temp.push({
              id: current.id,
              type: "bizDeploy",
              data: {
                name: current.name,
                config: current.config,
              },
            });
            break;

          case WorkflowNodeType.Notify:
            temp.push({
              id: current.id,
              type: "bizNotify",
              data: {
                name: current.name,
                config: current.config,
              },
            });
            break;

          case WorkflowNodeType.ExecuteResultBranch: {
            const tryNode = temp.pop()!;
            temp.push({
              id: current.id,
              type: "tryCatch",
              blocks: [
                {
                  id: current.branches?.find((b) => b.type === WorkflowNodeType.ExecuteSuccess)?.id || nanoid(),
                  type: "tryBlock",
                  blocks: [tryNode],
                  data: {
                    name: current.branches?.find((b) => b.type === WorkflowNodeType.ExecuteSuccess)?.name,
                  },
                },
                {
                  id: current.branches?.find((b) => b.type === WorkflowNodeType.ExecuteFailure)?.id || nanoid(),
                  type: "catchBlock",
                  blocks: [
                    ...convert(current.branches?.find((b) => b.type === WorkflowNodeType.ExecuteFailure)?.next),
                    {
                      id: nanoid(),
                      type: "end",
                      data: {
                        name: t("workflow_node.end.default_name"),
                      },
                    },
                  ],
                  data: {
                    name: current.branches?.find((b) => b.type === WorkflowNodeType.ExecuteFailure)?.name,
                  },
                },
              ],
              data: {
                name: current.name,
                config: current.config,
              },
            });

            current = current.branches?.find((b) => b.type === WorkflowNodeType.ExecuteSuccess);
            break;
          }

          case WorkflowNodeType.Branch: {
            temp.push({
              id: current.id,
              type: "condition",
              blocks:
                current.branches?.map((branch) => {
                  return {
                    id: branch.id,
                    type: "branchBlock",
                    blocks: convert(branch.next),
                    data: {
                      name: branch.name,
                      config: branch.config,
                    },
                  };
                }) ?? [],
              data: {
                name: current.name,
                config: current.config,
              },
            });
            break;
          }
        }

        current = current?.next;
      }

      return temp;
    };

    res.push(...convert(root));
  }

  res.push({
    id: nanoid(),
    type: "end",
    data: {
      name: t("workflow_node.end.default_name"),
    },
  });

  return res;
};

export default WorkflowDetailDesign;
