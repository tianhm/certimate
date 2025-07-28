import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { Outlet, useLocation, useNavigate, useParams } from "react-router-dom";
import { IconEdit, IconHistory, IconRobot } from "@tabler/icons-react";
import { App, Button, Input, type InputRef, Segmented, Skeleton } from "antd";

import Show from "@/components/Show";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";
import { mergeCls } from "@/utils/css";
import { getErrMsg } from "@/utils/error";

const WorkflowDetail = () => {
  const location = useLocation();
  const navigate = useNavigate();

  const { t } = useTranslation();

  const { id: workflowId } = useParams();
  const { workflow, initialized, ...workflowState } = useWorkflowStore(useZustandShallowSelector(["workflow", "initialized", "init", "destroy"]));
  useEffect(() => {
    workflowState.init(workflowId!);

    return () => {
      workflowState.destroy();
    };
  }, [workflowId]);

  const tabs = [
    ["design", "workflow.detail.design.tab", <IconRobot size="1em" />],
    ["runs", "workflow.detail.runs.tab", <IconHistory size="1em" />],
  ] satisfies [string, string, React.ReactElement][];
  const [tabValue, setTabValue] = useState<string>();
  useEffect(() => {
    const subpath = location.pathname.split("/")[3];
    if (!subpath) {
      navigate(`/workflows/${workflowId}/${tabs[0][0]}`, { replace: true });
      return;
    }

    setTabValue(subpath);
  }, [location.pathname, workflowId]);

  const handleTabChange = (value: string) => {
    setTabValue(value);
    navigate(`/workflows/${workflowId}/${value}`);

    workflowState.init(workflow.id); // reload state
  };

  return (
    <div className="flex size-full flex-col">
      <div className="px-6 py-4">
        <div className="relative mx-auto max-w-320">
          <WorkflowDetailBaseName />
          <WorkflowDetailBaseDescription />

          <div className="absolute -bottom-12 left-1/2 z-1 -translate-x-1/2">
            <Segmented
              className="shadow"
              options={tabs.map(([key, label, icon]) => ({
                value: key,
                label: <span className="px-2 text-sm">{t(label)}</span>,
                icon: (
                  <span className="anticon scale-125" role="img">
                    {icon}
                  </span>
                ),
              }))}
              size="large"
              value={tabValue}
              defaultValue="design"
              onChange={handleTabChange}
            />
          </div>
        </div>
      </div>

      <div className="p-4">
        <Show
          when={initialized}
          fallback={
            <div className="mx-auto max-w-320 pt-12">
              <Skeleton active />
            </div>
          }
        >
          <Outlet />
        </Show>
      </div>
    </div>
  );
};

const WorkflowDetailBaseName = () => {
  const { t } = useTranslation();

  const { notification } = App.useApp();

  const { workflow, setBaseInfo: setWorkflowBaseInfo } = useWorkflowStore(useZustandShallowSelector(["workflow", "setBaseInfo"]));

  const inputRef = useRef<InputRef>(null);

  const [editing, setEditing] = useState(false);
  const [value, setValue] = useState("");

  useEffect(() => {
    setEditing(false);
  }, [workflow.id]);

  const handleEditClick = () => {
    setEditing(true);
    setValue(workflow.name);
    setTimeout(() => {
      inputRef.current?.focus({ cursor: "all" });
    }, 0);
  };

  const handleValueChange = (value: string) => {
    setValue(value);
  };

  const handleValueConfirm = async (value: string) => {
    value = value.trim();
    if (value === (workflow.name || "")) {
      setEditing(false);
      return;
    }

    setEditing(false);

    try {
      await setWorkflowBaseInfo(value, workflow.description!);
    } catch (err) {
      notification.error({ message: t("common.text.request_error"), description: getErrMsg(err) });

      throw err;
    }
  };

  return (
    <div className="group relative flex items-center gap-1">
      <h1
        className={mergeCls({
          invisible: editing,
        })}
      >
        {workflow.name || "\u00A0"}
      </h1>
      <Button
        className={mergeCls("mb-2 opacity-0 transition-opacity group-hover:opacity-100", {
          invisible: editing,
        })}
        icon={<IconEdit size="1.25em" stroke="1.25" />}
        type="text"
        onClick={handleEditClick}
      />
      <Input
        className={mergeCls("absolute top-0 left-0", editing ? "block" : "hidden")}
        ref={inputRef}
        size="large"
        value={value}
        variant="filled"
        onBlur={(e) => handleValueConfirm(e.target.value)}
        onChange={(e) => handleValueChange(e.target.value)}
        onPressEnter={(e) => e.currentTarget.blur()}
      />
    </div>
  );
};

const WorkflowDetailBaseDescription = () => {
  const { t } = useTranslation();

  const { notification } = App.useApp();

  const { workflow, setBaseInfo: setWorkflowBaseInfo } = useWorkflowStore(useZustandShallowSelector(["workflow", "setBaseInfo"]));

  const inputRef = useRef<InputRef>(null);

  const [editing, setEditing] = useState(false);
  const [value, setValue] = useState("");

  useEffect(() => {
    setEditing(false);
  }, [workflow.id]);

  const handleEditClick = () => {
    setEditing(true);
    setValue(workflow.description || "");
    setTimeout(() => {
      inputRef.current?.focus({ cursor: "all" });
    }, 0);
  };

  const handleValueChange = (value: string) => {
    setValue(value);
  };

  const handleValueConfirm = async (value: string) => {
    value = value.trim();
    if (value === (workflow.description || "")) {
      setEditing(false);
      return;
    }

    setEditing(false);

    try {
      await setWorkflowBaseInfo(workflow.name!, value);
    } catch (err) {
      notification.error({ message: t("common.text.request_error"), description: getErrMsg(err) });

      throw err;
    }
  };

  return (
    <div className="group relative flex items-center gap-1">
      <p
        className={mergeCls("text-base text-gray-500", {
          invisible: editing,
        })}
      >
        {workflow.description || "\u00A0"}
      </p>
      <Button
        className={mergeCls("mb-4 opacity-0 transition-opacity group-hover:opacity-100", {
          invisible: editing,
        })}
        icon={<IconEdit size="1.25em" stroke="1.25" />}
        type="text"
        onClick={handleEditClick}
      />
      <Input
        className={mergeCls("absolute top-0 left-0", editing ? "block" : "hidden")}
        ref={inputRef}
        value={value}
        variant="filled"
        onBlur={(e) => handleValueConfirm(e.target.value)}
        onChange={(e) => handleValueChange(e.target.value)}
        onPressEnter={(e) => e.currentTarget.blur()}
      />
    </div>
  );
};

export default WorkflowDetail;
