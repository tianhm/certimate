import { createContext } from "react";
import { type NodeRenderReturnType } from "@flowgram.ai/fixed-layout-editor";

export const NodeRenderContext = createContext<NodeRenderReturnType>({} as NodeRenderReturnType);
