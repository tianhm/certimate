import { useEffect } from "react";
import { type NodeRenderProps, useClientContext, useNodeRender, useRefresh } from "@flowgram.ai/fixed-layout-editor";

import { NodeRenderContext } from "./NodeRenderContext";
import { type NodeRegistry } from "./nodes/typings";

export interface NodeProps extends NodeRenderProps {}

const Node = (_: NodeProps) => {
  const ctx = useClientContext();

  const nodeRender = useNodeRender();

  const refresh = useRefresh();
  useEffect(() => {
    const disposable = nodeRender.form?.onFormValuesChange?.(() => refresh());
    return () => disposable?.dispose();
  }, [nodeRender.form]);
  useEffect(() => {
    const toDispose = ctx.document.originTree.onTreeChange(() => refresh());
    return () => toDispose.dispose();
  }, []);

  return (
    <div
      style={{
        opacity: nodeRender.dragging ? 0.3 : 1,
        outline: nodeRender.form?.state?.invalid ? "1px solid var(--color-error)" : "none",
        ...nodeRender.node.getNodeRegistry<NodeRegistry>().meta?.style,
      }}
      onMouseEnter={nodeRender.onMouseEnter}
      onMouseLeave={nodeRender.onMouseLeave}
      onMouseDown={(e) => {
        nodeRender.startDrag(e);
        e.stopPropagation();
      }}
    >
      <NodeRenderContext.Provider value={nodeRender}>{nodeRender.form?.render()}</NodeRenderContext.Provider>
    </div>
  );
};

export default Node;
