import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { IconArrowBackUp, IconDots } from "@tabler/icons-react";
import { Alert, App, Button, Card, Dropdown, Space } from "antd";
import { isEqual } from "radash";

import Show from "@/components/Show";
import WorkflowElementsContainer from "@/components/workflow/WorkflowElementsContainer";
import { isAllNodesValidated } from "@/domain/workflow";
import { WORKFLOW_RUN_STATUSES } from "@/domain/workflowRun";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";
import { getErrMsg } from "@/utils/error";

/**
 *
 * @deprecated
 */
const WorkflowDetailDesign = () => {
  const { t } = useTranslation();

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

  const handleRollbackClick = () => {
    modal.confirm({
      title: t("workflow.detail.design.action.rollback.modal.title"),
      content: t("workflow.detail.design.action.rollback.modal.content"),
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
      message.warning(t("workflow.detail.design.uncompleted_design.alert"));
      return;
    }

    modal.confirm({
      title: t("workflow.detail.design.action.publish.modal.title"),
      content: t("workflow.detail.design.action.publish.modal.content"),
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
    <div className="min-h-[360px] flex-1 overflow-hidden">
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
        <div className="size-full pt-9">
          <div className="absolute top-8 z-10 w-full px-4">
            <div className="container">
              <Alert
                className="mb-2"
                message={
                  <div>
                    该子页面即将在 v0.4.0 中被移除。
                    <br />
                    This subpage will be dropped in v0.4.0.
                  </div>
                }
                showIcon
                closable
                type="warning"
              />
              <div className="flex items-center justify-end gap-4">
                <div className="flex flex-1 items-center justify-end gap-4 overflow-hidden">
                  <div className="flex-1 overflow-hidden">
                    <Show when={workflow.hasDraft!}>
                      <Alert message={<div className="truncate">{t("workflow.detail.design.unpublished_draft.alert")}</div>} showIcon type="warning" />
                    </Show>
                  </div>
                  <Space.Compact>
                    <Button disabled={!allowPublish} type="primary" onClick={handlePublishClick}>
                      {t("workflow.detail.design.action.publish.button")}
                    </Button>
                    <Dropdown
                      menu={{
                        items: [
                          {
                            key: "rollback",
                            disabled: !allowRollback,
                            label: t("workflow.detail.design.action.rollback.button"),
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

          <WorkflowElementsContainer className="pt-12" />
        </div>
      </Card>
    </div>
  );
};

export default WorkflowDetailDesign;
