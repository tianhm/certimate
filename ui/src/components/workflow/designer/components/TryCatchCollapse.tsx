import { useState } from "react";
import {
  type CustomLabelProps,
  FlowNodeRenderData,
  FlowNodeTransformData,
  FlowRendererRegistry,
  FlowTextKey,
  useBaseColor,
} from "@flowgram.ai/fixed-layout-editor";

export interface TryCatchCollapseProps extends CustomLabelProps {}

const TryCatchCollapse = ({ node, ...props }: TryCatchCollapseProps) => {
  const { baseColor, baseActivatedColor } = useBaseColor();

  const nodeActivateData = node.getData(FlowNodeRenderData)!;
  const nodeTransformData = node.getData(FlowNodeTransformData)!;

  const [hoverActivated, setHoverActivated] = useState(false);

  const handleMouseEnter = () => {
    setHoverActivated(true);
    nodeActivateData.activated = true;
  };

  const handleMouseLeave = () => {
    setHoverActivated(false);
    nodeActivateData.activated = false;
  };

  if (!nodeTransformData || !nodeTransformData.parent) {
    return <></>;
  }

  const width = nodeTransformData.inputPoint.x - nodeTransformData.parent.inputPoint.x;
  const height = 40;
  return (
    <div
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      style={{
        width,
        height,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        gap: 6,
      }}
    >
      <div
        data-label-id={props.labelId}
        style={{
          fontSize: 12,
          color: hoverActivated ? baseActivatedColor : baseColor,
          textAlign: "center",
          lineHeight: "20px",
          whiteSpace: "nowrap",
          backgroundColor: "var(--g-editor-background)",
        }}
      >
        {node.getService(FlowRendererRegistry).getText(FlowTextKey.CATCH_TEXT)}
      </div>

      {/* {renderCollapse()} */}
    </div>
  );
};

export default TryCatchCollapse;
