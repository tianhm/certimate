import { forwardRef, useImperativeHandle, useMemo, useRef } from "react";
import {
  ConstantKeys,
  EditorRenderer,
  FixedLayoutEditorProvider,
  type FixedLayoutPluginContext,
  type FixedLayoutProps,
  type FlowDocumentJSON,
  type FlowNodeEntity,
  FlowTextKey,
  getNodeForm,
} from "@flowgram.ai/fixed-layout-editor";
import { createMinimapPlugin } from "@flowgram.ai/minimap-plugin";
import "@flowgram.ai/fixed-layout-editor/index.css";
import { theme } from "antd";

import { getFlowComponents } from "./components";
import { EditorContextProvider } from "./EditorContext";
import NodeRender from "./NodeRender";
import { getFlowNodeRegistries } from "./nodes";
import { BranchNode } from "./nodes/_shared";
import "./flowgram.css";

export interface EditorProps {
  className?: string;
  style?: React.CSSProperties;
  children?: React.ReactNode;
  initialData?: FlowDocumentJSON;
  readonly?: boolean;
  onNodeClick?: (ctx: FixedLayoutPluginContext, node: FlowNodeEntity) => void;
}

export interface EditorInstance extends FixedLayoutPluginContext {
  validateAllNodes(): Promise<boolean>;
}

const Editor = forwardRef<EditorInstance, EditorProps>(({ className, style, children, initialData, readonly, onNodeClick }, ref) => {
  const { token: themeToken } = theme.useToken();

  const flowgramEditorRef = useRef<FixedLayoutPluginContext>(null);
  const flowgramEditorProps = useMemo<FixedLayoutProps>(
    () => ({
      initialData: initialData,

      constants: {
        [ConstantKeys.BASE_COLOR]: themeToken.colorBorder,
        [ConstantKeys.BASE_ACTIVATED_COLOR]: themeToken.colorPrimary,
        [ConstantKeys.NODE_SPACING]: 48,
        [ConstantKeys.BRANCH_SPACING]: 48,
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

      selectBox: {
        enable: false,
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
            render: () => <BranchNode>{type}</BranchNode>,
          },
        };
      },

      plugins: () => [
        createMinimapPlugin({
          disableLayer: true,
          enableDisplayAllNodes: true,
          canvasStyle: {
            canvasWidth: 160,
            canvasHeight: 160,
          },
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

  useImperativeHandle(ref, () => {
    return {
      get container() {
        return flowgramEditorRef.current!.container;
      },
      get document() {
        return flowgramEditorRef.current!.document;
      },
      get playground() {
        return flowgramEditorRef.current!.playground;
      },
      get operation() {
        return flowgramEditorRef.current!.operation;
      },
      get clipboard() {
        return flowgramEditorRef.current!.clipboard;
      },
      get selection() {
        return flowgramEditorRef.current!.selection;
      },
      get history() {
        return flowgramEditorRef.current!.history;
      },

      get(identifier) {
        return flowgramEditorRef.current!.get(identifier);
      },
      getAll(identifier) {
        return flowgramEditorRef.current!.getAll(identifier);
      },
      validateAllNodes() {
        const nodes = flowgramEditorRef.current!.document.getAllNodes();
        const forms = nodes.map((node) => getNodeForm(node)).filter((form) => form != null);
        return Promise.allSettled(forms.map((form) => form.validate())).then((res) => forms.every((form, index) => res[index] && !form.state.invalid));
      },
    };
  });

  return (
    <FixedLayoutEditorProvider ref={flowgramEditorRef} {...flowgramEditorProps}>
      <EditorContextProvider value={{ onNodeClick: (node) => onNodeClick?.(flowgramEditorRef.current!, node) }}>
        <EditorRenderer className={className} style={style} />
        {children}
      </EditorContextProvider>
    </FixedLayoutEditorProvider>
  );
});

export default Editor;
