import { type WorkflowModel } from "@/domain/workflow";

export interface WorkflowState {
  workflow: WorkflowModel;
  initialized: boolean;
}

export interface WorkflowActions {
  init(id: string): void;
  destroy(): void;

  setName: (name: Required<WorkflowModel>["name"]) => void;
  setDescription: (description: Required<WorkflowModel>["description"]) => void;
  setEnabled(enabled: Required<WorkflowModel>["enabled"]): void;
  setDraft(draft: Required<WorkflowModel>["draft"]): void;

  publish(): void;
  rollback(): void;
}

export interface WorkflowStore extends WorkflowState, WorkflowActions {}
