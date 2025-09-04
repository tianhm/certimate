import { useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useMount } from "ahooks";
import { Avatar, Card, Checkbox, Empty, Flex, Input, type InputRef, Tooltip, Typography } from "antd";

import Show from "@/components/Show";
import { type NotificationProvider, notificationProvidersMap } from "@/domain/provider";
import { useZustandShallowSelector } from "@/hooks";
import { useAccessesStore } from "@/stores/access";
import { mergeCls } from "@/utils/css";

import { type SharedPickerProps, usePickerWrapperCols } from "./_shared";

export interface NotificationProviderPickerProps extends SharedPickerProps<NotificationProvider> {}

const NotificationProviderPicker = ({
  className,
  style,
  autoFocus,
  gap = "middle",
  placeholder,
  showSearch = true,
  onFilter,
  onSelect,
}: NotificationProviderPickerProps) => {
  const { t } = useTranslation();

  const { accesses, fetchAccesses } = useAccessesStore(useZustandShallowSelector(["accesses", "fetchAccesses"]));
  useMount(() => fetchAccesses(false));

  const { wrapperElRef, cols } = usePickerWrapperCols(320);

  const [isAvailableOnly, setIsAvailableOnly] = useState(true);

  const [keyword, setKeyword] = useState<string>();
  const keywordInputRef = useRef<InputRef>(null);
  useMount(() => {
    if (autoFocus) {
      setTimeout(() => keywordInputRef.current?.focus(), 1);
    }
  });

  const providers = useMemo(() => {
    return Array.from(notificationProvidersMap.values())
      .filter((provider) => {
        if (onFilter) {
          return onFilter(provider.type, provider);
        }

        return true;
      })
      .filter((provider) => {
        if (isAvailableOnly) {
          return accesses.some((access) => access.provider === provider.provider);
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
  }, [onFilter, accesses, isAvailableOnly, keyword]);

  const handleProviderTypeSelect = (value: string) => {
    onSelect?.(value);
  };

  return (
    <div className={className} style={style} ref={wrapperElRef}>
      <Show when={showSearch}>
        <div className="mb-4">
          <Input.Search ref={keywordInputRef} placeholder={placeholder ?? t("common.text.search")} onChange={(e) => setKeyword(e.target.value.trim())} />
        </div>
      </Show>

      <Flex className="mb-4" justify="end">
        <Checkbox checked={isAvailableOnly} onClick={() => setIsAvailableOnly(!isAvailableOnly)}>
          {t("provider.text.show_available_provider_only")}
        </Checkbox>
      </Flex>

      <Show when={providers.length > 0} fallback={<Empty description={t("provider.text.nodata")} image={Empty.PRESENTED_IMAGE_SIMPLE} />}>
        <div
          className={mergeCls("grid w-full gap-2", `grid-cols-${cols}`, {
            "gap-4": gap === "large",
            "gap-2": gap === "middle",
            "gap-1": gap === "small",
            [`gap-${+gap || "2"}`]: typeof gap === "number",
          })}
        >
          {providers.map((provider) => {
            return (
              <div key={provider.type}>
                <Card
                  className="h-16 w-full overflow-hidden shadow"
                  styles={{ body: { height: "100%", padding: "0.5rem 1rem" } }}
                  hoverable
                  onClick={() => {
                    handleProviderTypeSelect(provider.type);
                  }}
                >
                  <div className="flex size-full items-center gap-4 overflow-hidden">
                    <Avatar className="bg-stone-100" icon={<img src={provider.icon} />} shape="square" size={28} />
                    <div className="flex-1 overflow-hidden">
                      <div className="line-clamp-2 max-w-full">
                        <Tooltip title={t(provider.name)} mouseEnterDelay={1}>
                          <Typography.Text>{t(provider.name) || "\u00A0"}</Typography.Text>
                        </Tooltip>
                      </div>
                    </div>
                  </div>
                </Card>
              </div>
            );
          })}
        </div>
      </Show>
    </div>
  );
};

export default NotificationProviderPicker;
