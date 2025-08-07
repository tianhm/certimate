import { useEffect } from "react";
import { type NodeRenderProps, useClientContext, useNodeRender, useRefresh } from "@flowgram.ai/fixed-layout-editor";

import { useEditorContext } from "./EditorContext";
import { NodeRenderContextProvider } from "./NodeRenderContext";
import { type NodeRegistry } from "./nodes/typings";

export interface NodeProps extends NodeRenderProps {}

const Node = (_: NodeProps) => {
  const ctx = useClientContext();

  const refresh = useRefresh();

  const nodeRender = useNodeRender();

  useEffect(() => {
    const d = ctx.document.originTree.onTreeChange(() => refresh());
    return () => d.dispose();
  }, []);

  useEffect(() => {
    const d1 = nodeRender.form?.onFormValuesChange?.(() => refresh());
    const d2 = nodeRender.form?.onValidate?.(() => refresh());
    return () => {
      d1?.dispose();
      d2?.dispose();
    };
  }, [nodeRender.form]);

  const { onNodeClick } = useEditorContext();

  return (
    <div
      style={{
        opacity: nodeRender.dragging ? 0.3 : 1,
        ...nodeRender.node.getNodeRegistry<NodeRegistry>().meta?.style,
      }}
      onMouseEnter={nodeRender.onMouseEnter}
      onMouseLeave={nodeRender.onMouseLeave}
      onMouseDown={(e) => {
        nodeRender.startDrag(e);
        e.stopPropagation();
      }}
      onClick={() => {
        onNodeClick?.(nodeRender.node);
      }}
    >
      <NodeRenderContextProvider value={nodeRender}>{nodeRender.form?.render()}</NodeRenderContextProvider>
    </div>
  );
};

export default Node;
