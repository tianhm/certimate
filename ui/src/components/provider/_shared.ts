import { useMemo, useRef } from "react";
import { useSize } from "ahooks";
import { type SelectProps } from "antd";

type Provider = { type: string; name: string };

export interface SharedSelectProps<T extends Provider>
  extends Omit<SelectProps, "filterOption" | "filterSort" | "labelRender" | "options" | "optionFilterProp" | "optionLabelProp" | "optionRender"> {
  className?: string;
  style?: React.CSSProperties;
  onFilter?: (value: string, option: T) => boolean;
}

export interface SharedPickerProps<T extends Provider> {
  className?: string;
  style?: React.CSSProperties;
  autoFocus?: boolean;
  gap?: number | "small" | "middle" | "large";
  placeholder?: string;
  showSearch?: boolean;
  onFilter?: (value: string, option: T) => boolean;
  onSelect?: (value: string) => void;
}

export const usePickerWrapperCols = (width: number) => {
  const wrapperElRef = useRef<HTMLDivElement>(null);
  const wrapperSize = useSize(wrapperElRef);

  const cols = useMemo(() => {
    if (!wrapperSize) {
      return 1;
    }

    const cols = Math.floor(wrapperSize.width / width);
    return Math.min(9, Math.max(1, cols));
  }, [wrapperSize, width]);

  return {
    wrapperElRef,
    cols,
  };
};
