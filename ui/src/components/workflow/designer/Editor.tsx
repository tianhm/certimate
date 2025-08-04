import { useMemo, useRef } from "react";
import {
  ConstantKeys,
  EditorRenderer,
  FixedLayoutEditorProvider,
  type FixedLayoutPluginContext,
  type FixedLayoutProps,
  type FlowDocumentJSON,
  FlowTextKey,
} from "@flowgram.ai/fixed-layout-editor";
import { createMinimapPlugin } from "@flowgram.ai/minimap-plugin";
import "@flowgram.ai/fixed-layout-editor/index.css";
import { useDeepCompareEffect } from "ahooks";
import { theme } from "antd";

import { getFlowComponents } from "./components";
import NodeRender from "./NodeRender";
import { getFlowNodeRegistries } from "./nodes";
import { BlockNode } from "./nodes/_shared";
import "./flowgram.css";

export interface EditorProps {
  className?: string;
  style?: React.CSSProperties;
  children?: React.ReactNode;
  initialData?: FlowDocumentJSON;
  readonly?: boolean;
}

const Editor = ({ className, style, children, initialData, readonly }: EditorProps) => {
  const { token: themeToken } = theme.useToken();

  const flowgramEditorRef = useRef<FixedLayoutPluginContext>(null);
  const flowgramEditorProps = useMemo<FixedLayoutProps>(
    () => ({
      initialData: initialData,

      constants: {
        [ConstantKeys.BASE_COLOR]: themeToken.colorBorder,
        [ConstantKeys.BASE_ACTIVATED_COLOR]: themeToken.colorPrimary,
        [ConstantKeys.NODE_SPACING]: 64,
        [ConstantKeys.BRANCH_SPACING]: 64,
        // [ConstantKeys.INLINE_BLOCKS_PADDING_TOP]: 48,
        // [ConstantKeys.INLINE_BLOCKS_PADDING_BOTTOM]: 48,
      },

      background: {
        backgroundColor: themeToken.colorBgContainer,
        dotSize: 0,
      },

      playground: {
        autoFocus: true,
        autoResize: true,
        preventGlobalGesture: true,
      },

      scroll: {
        enableScrollLimit: true,
      },

      readonly: readonly,

      nodeEngine: {
        enable: true,
      },

      variableEngine: {
        enable: true,
      },

      materials: {
        components: getFlowComponents(),
        renderTexts: {
          [FlowTextKey.TRY_START_TEXT]: "Try",
          [FlowTextKey.TRY_END_TEXT]: "Finally",
          [FlowTextKey.CATCH_TEXT]: "Catch",
        },
        renderDefaultNode: NodeRender,
      },

      nodeRegistries: getFlowNodeRegistries(),

      getNodeDefaultRegistry(type) {
        return {
          type,
          meta: {
            defaultExpanded: true,
          },
          formMeta: {
            render: () => <BlockNode>{type}</BlockNode>,
          },
        };
      },

      plugins: () => [
        createMinimapPlugin({
          disableLayer: true,
          enableDisplayAllNodes: true,
        }),
      ],

      onAllLayersRendered: (ctx) => {
        // 画布初始化后向下滚动一点，露出可能被 Alert 遮挡的部分
        setTimeout(() => {
          ctx.playground.config.scroll({ scrollY: -80 });
        }, 1);
      },
    }),
    [themeToken, initialData, readonly]
  );

  useDeepCompareEffect(() => {
    flowgramEditorRef.current?.document?.fromJSON?.(initialData ?? {});
  }, [initialData]);

  return (
    <div className={className} style={style}>
      <FixedLayoutEditorProvider ref={flowgramEditorRef} {...flowgramEditorProps}>
        <EditorRenderer />
        {children}
      </FixedLayoutEditorProvider>
    </div>
  );
};

export default Editor;
