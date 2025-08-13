import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { EditorState, FlowLayoutDefault, useClientContext, usePlaygroundTools, useRefresh } from "@flowgram.ai/fixed-layout-editor";
import { IconHandStop, IconLayoutCards, IconMatrix, IconMaximize, IconMinus, IconPlus } from "@tabler/icons-react";
import { Button, Dropdown, Tooltip } from "antd";

import { mergeCls } from "@/utils/css";

import Minimap from "./Minimap";

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
    const d = playground.config.onReadonlyOrDisabledChange(() => refresh());

    return () => d.dispose();
  }, [playground]);

  const [isMinimapVisible, setIsMinimapVisible] = useState(window.screen.availWidth >= 1024);

  const [isMouseFriendly, setIsMouseFriendly] = useState(playground.editorState.is(EditorState.STATE_MOUSE_FRIENDLY_SELECT.id));

  const handleToggleLayout = useCallback(() => {
    if (tools.isVertical) {
      tools.changeLayout(FlowLayoutDefault.HORIZONTAL_FIXED_LAYOUT);
    } else {
      tools.changeLayout(FlowLayoutDefault.VERTICAL_FIXED_LAYOUT);
    }
  }, [tools.isVertical]);

  const handleToggleMinimap = useCallback(() => {
    setIsMinimapVisible((prev) => !prev);
  }, [isMinimapVisible]);

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
      <div className="relative flex items-center gap-2">
        <Tooltip title={isMouseFriendly ? t("workflow.detail.design.toolbar.hand_mode") : t("workflow.detail.design.toolbar.pointer_mode")}>
          <Button
            ghost={isMouseFriendly}
            icon={<IconHandStop size="1.25em" />}
            type={isMouseFriendly ? "primary" : "default"}
            onClick={handleToggleMouseFriendly}
          />
        </Tooltip>

        <Tooltip title={t("workflow.detail.design.toolbar.zoomout")}>
          <Button icon={<IconMinus size="1.25em" />} onClick={() => tools.zoomout()} />
        </Tooltip>
        <Dropdown
          menu={{
            items: [
              ...[200, 100, 75, 50, 25].map((zoom) => ({
                key: `${zoom}%`,
                label: `${zoom}%`,
                onClick: () => tools.updateZoom(zoom / 100),
              })),
              {
                type: "divider",
              },
              {
                key: "auto",
                label: t("workflow.detail.design.toolbar.auto_fit"),
                onClick: () => tools.fitView(),
              },
            ],
          }}
          trigger={["click"]}
        >
          <Button className="w-16 text-center">{Math.round(tools.zoom * 100)}%</Button>
        </Dropdown>
        <Tooltip title={t("workflow.detail.design.toolbar.zoomin")}>
          <Button icon={<IconPlus size="1.25em" />} onClick={() => tools.zoomin()} />
        </Tooltip>
        <Tooltip title={t("workflow.detail.design.toolbar.auto_fit")}>
          <Button icon={<IconMaximize size="1.25em" />} onClick={() => tools.fitView()} />
        </Tooltip>

        <Tooltip title={tools.isVertical ? t("workflow.detail.design.toolbar.vertical_layout") : t("workflow.detail.design.toolbar.horizontal_layout")}>
          <Button icon={<IconLayoutCards className={mergeCls({ ["rotate-90"]: tools.isVertical })} size="1.25em" />} onClick={handleToggleLayout} />
        </Tooltip>

        <Tooltip title={t("workflow.detail.design.toolbar.minimap")}>
          <Button icon={<IconMatrix size="1.25em" />} ghost={isMinimapVisible} type={isMinimapVisible ? "primary" : "default"} onClick={handleToggleMinimap} />
        </Tooltip>
        {isMinimapVisible && <Minimap className="absolute right-0 bottom-[42px]" />}
      </div>
    </div>
  );
};

export default Toolbar;
