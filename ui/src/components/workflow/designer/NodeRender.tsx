import { type NodeRenderProps, useNodeRender } from "@flowgram.ai/fixed-layout-editor";

export interface NodeProps extends NodeRenderProps {}

const Node = (_: NodeProps) => {
  const nodeRender = useNodeRender();

  return (
    <div
      style={{
        opacity: nodeRender.dragging ? 0.3 : 1,
        outline: nodeRender.form?.state?.invalid ? "1px solid var(--color-error)" : "none",
        ...nodeRender.node.getNodeRegistry().meta.style,
      }}
      onMouseEnter={nodeRender.onMouseEnter}
      onMouseLeave={nodeRender.onMouseLeave}
      onMouseDown={(e) => {
        nodeRender.startDrag(e);
        e.stopPropagation();
      }}
    >
      {nodeRender.form?.render()}
    </div>
  );
};

export default Node;
