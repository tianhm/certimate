import { createContext, useContext } from "react";
import { type FlowNodeEntity } from "@flowgram.ai/fixed-layout-editor";

export type NodeRenderContextType = {
  onNodeClick: (node: FlowNodeEntity) => void;
};

export const EditorContext = createContext<NodeRenderContextType>({
  onNodeClick: () => {},
});

export const EditorContextProvider = EditorContext.Provider;

export const useEditorContext = () => {
  const context = useContext(EditorContext);
  if (!context) {
    throw new Error("`EditorContext` must be used within a `EditorContextProvider`");
  }
  return context;
};
