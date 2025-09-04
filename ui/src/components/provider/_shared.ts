import { useMemo, useRef } from "react";
import { useTranslation } from "react-i18next";
import { useMount, useSize } from "ahooks";

import { type SelectProps } from "antd";

import { useZustandShallowSelector } from "@/hooks";
import { useAccessesStore } from "@/stores/access";

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

export const usePickerDataSource = <T extends Provider>({
  dataSource,
  filters,
  keyword,
  onFilter,
  deps,
}: {
  dataSource: T[];
  filters?: Array<(option: T) => boolean>;
  keyword?: string;
  onFilter?: (value: string, option: T) => boolean;
  deps?: React.DependencyList;
}) => {
  const { t } = useTranslation();

  const { accesses, fetchAccesses } = useAccessesStore(useZustandShallowSelector(["accesses", "fetchAccesses"]));
  useMount(() => fetchAccesses(false));

  const filteredDataSource = useMemo(() => {
    return dataSource
      .filter((provider) => {
        if (onFilter) {
          return onFilter(provider.type, provider);
        }

        return true;
      })
      .filter((provider) => {
        if (filters) {
          for (const filter of filters) {
            if (!filter(provider)) return false;
          }
        }

        return true;
      })
      .filter((provider) => {
        if (keyword) {
          const value = keyword.toLowerCase();
          return provider.type.toLowerCase().includes(value) || t(provider.name).toLowerCase().includes(value);
        }

        return true;
      });
  }, [dataSource, filters, keyword, onFilter, ...(deps ?? [])]);

  const availableDataSource = useMemo(() => {
    return filteredDataSource.filter((provider) => {
      return accesses.some((access) => {
        if ("builtin" in provider && provider.builtin) return true;
        if ("provider" in provider) return access.provider === provider.provider;
        return access.provider === provider.type;
      });
    });
  }, [accesses, filteredDataSource, ...(deps ?? [])]);

  const unavailableDataSource = useMemo(() => {
    return filteredDataSource.filter((item) => !availableDataSource.includes(item));
  }, [filteredDataSource, availableDataSource, ...(deps ?? [])]);

  return {
    all: dataSource,
    filtered: filteredDataSource,
    available: availableDataSource,
    unavailable: unavailableDataSource,
  };
};
