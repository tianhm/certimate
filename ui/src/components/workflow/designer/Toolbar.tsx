import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { EditorState, FlowLayoutDefault, useClientContext, usePlaygroundTools, useRefresh } from "@flowgram.ai/fixed-layout-editor";
import { IconHandStop, IconLayoutCards, IconMaximize, IconMinus, IconPlus, IconPointer } from "@tabler/icons-react";
import { Button, Tooltip } from "antd";

import { mergeCls } from "@/utils/css";

export interface ToolbarProps {
  className?: string;
  style?: React.CSSProperties;
}

const Toolbar = ({ className, style }: ToolbarProps) => {
  const { t } = useTranslation();

  const ctx = useClientContext();
  const { playground } = ctx;

  const tools = usePlaygroundTools({ minZoom: 0.1, maxZoom: 3, padding: 48 });

  const refresh = useRefresh();
  useEffect(() => {
    const disposable = playground.config.onReadonlyOrDisabledChange(() => refresh());
    return () => disposable.dispose();
  }, [playground]);

  const [isMouseFriendly, setIsMouseFriendly] = useState(playground.editorState.is(EditorState.STATE_MOUSE_FRIENDLY_SELECT.id));

  const handleToggleLayout = useCallback(() => {
    if (tools.isVertical) {
      tools.changeLayout(FlowLayoutDefault.HORIZONTAL_FIXED_LAYOUT);
    } else {
      tools.changeLayout(FlowLayoutDefault.VERTICAL_FIXED_LAYOUT);
    }
  }, [tools.isVertical]);

  const handleToggleMouseFriendly = useCallback(() => {
    if (isMouseFriendly) {
      playground.editorState.changeState(EditorState.STATE_SELECT.id);
      setIsMouseFriendly(false);
    } else {
      playground.editorState.changeState(EditorState.STATE_MOUSE_FRIENDLY_SELECT.id);
      setIsMouseFriendly(true);
    }
  }, [isMouseFriendly]);

  return (
    <div className={className} style={style}>
      <div className="flex items-center gap-2">
        <Tooltip title={t("workflow.detail.design.toolbar.zoomout")}>
          <Button icon={<IconMinus size="1.25em" />} onClick={() => tools.zoomout()} />
        </Tooltip>
        <Tooltip title={t("workflow.detail.design.toolbar.zoom")}>
          <Button className="w-16 text-center" onClick={() => tools.updateZoom(1)}>
            {Math.round(tools.zoom * 100)}%
          </Button>
        </Tooltip>
        <Tooltip title={t("workflow.detail.design.toolbar.zoomin")}>
          <Button icon={<IconPlus size="1.25em" />} onClick={() => tools.zoomin()} />
        </Tooltip>
        <Tooltip title={t("workflow.detail.design.toolbar.auto_fit")}>
          <Button icon={<IconMaximize size="1.25em" />} onClick={() => tools.fitView()} />
        </Tooltip>

        <Tooltip title={tools.isVertical ? t("workflow.detail.design.toolbar.vertical_layout") : t("workflow.detail.design.toolbar.horizontal_layout")}>
          <Button icon={<IconLayoutCards className={mergeCls({ ["rotate-90"]: tools.isVertical })} size="1.25em" />} onClick={handleToggleLayout} />
        </Tooltip>

        <Tooltip title={isMouseFriendly ? t("workflow.detail.design.toolbar.hand_mode") : t("workflow.detail.design.toolbar.pointer_mode")}>
          <Button icon={isMouseFriendly ? <IconHandStop size="1.25em" /> : <IconPointer size="1.25em" />} onClick={handleToggleMouseFriendly} />
        </Tooltip>
      </div>
    </div>
  );
};

export default Toolbar;
