import { useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useMount, useSize } from "ahooks";
import { Avatar, Card, Empty, Input, type InputRef, Tag, Tooltip, Typography } from "antd";

import Show from "@/components/Show";
import { ACCESS_USAGES, type AccessProvider, type AccessUsageType, accessProvidersMap } from "@/domain/provider";
import { mergeCls } from "@/utils/css";

export interface AccessProviderPickerProps {
  className?: string;
  style?: React.CSSProperties;
  autoFocus?: boolean;
  gap?: number | "small" | "middle" | "large";
  placeholder?: string;
  showOptionTags?: boolean | { [key in AccessUsageType]?: boolean };
  onFilter?: (value: string, option: AccessProvider) => boolean;
  onSelect?: (value: string) => void;
}

const AccessProviderPicker = ({ className, style, autoFocus, gap = "middle", placeholder, showOptionTags, onFilter, onSelect }: AccessProviderPickerProps) => {
  const { t } = useTranslation();

  const wrapperRef = useRef<HTMLDivElement>(null);
  const wrapperSize = useSize(wrapperRef);

  const [keyword, setKeyword] = useState<string>();
  const keywordInputRef = useRef<InputRef>(null);
  useMount(() => {
    if (autoFocus) {
      setTimeout(() => keywordInputRef.current?.focus(), 1);
    }
  });

  const showOptionTagForDNS = useMemo(() => {
    return typeof showOptionTags === "object" ? !!showOptionTags?.[ACCESS_USAGES.DNS] : !!showOptionTags;
  }, [showOptionTags]);
  const showOptionTagForHosting = useMemo(() => {
    return typeof showOptionTags === "object" ? !!showOptionTags?.[ACCESS_USAGES.HOSTING] : !!showOptionTags;
  }, [showOptionTags]);
  const showOptionTagForCA = useMemo(() => {
    return typeof showOptionTags === "object" ? !!showOptionTags?.[ACCESS_USAGES.CA] : !!showOptionTags;
  }, [showOptionTags]);
  const showOptionTagForNotification = useMemo(() => {
    return typeof showOptionTags === "object" ? !!showOptionTags?.[ACCESS_USAGES.NOTIFICATION] : !!showOptionTags;
  }, [showOptionTags]);
  const showOptionTagAnyhow = useMemo(() => {
    return showOptionTagForDNS || showOptionTagForHosting || showOptionTagForCA || showOptionTagForNotification;
  }, [showOptionTagForDNS, showOptionTagForHosting, showOptionTagForCA, showOptionTagForNotification]);

  const providers = useMemo(() => {
    return Array.from(accessProvidersMap.values())
      .filter((provider) => {
        if (onFilter) {
          return onFilter(provider.type, provider);
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
  }, [onFilter, keyword]);
  const providerCols = useMemo(() => {
    if (!wrapperSize) {
      return 1;
    }

    const cols = Math.floor(wrapperSize.width / (showOptionTagAnyhow ? 240 : 200));
    return Math.min(9, Math.max(1, cols));
  }, [wrapperSize, showOptionTagAnyhow]);

  const handleProviderTypeSelect = (value: string) => {
    onSelect?.(value);
  };

  return (
    <div className={className} style={style} ref={wrapperRef}>
      <Input.Search ref={keywordInputRef} placeholder={placeholder ?? t("common.text.search")} onChange={(e) => setKeyword(e.target.value.trim())} />

      <div className="mt-4">
        <Show when={providers.length > 0} fallback={<Empty description={t("common.text.nodata")} image={Empty.PRESENTED_IMAGE_SIMPLE} />}>
          <div
            className={mergeCls("grid w-full gap-2", `grid-cols-${providerCols}`, {
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
                    className={mergeCls("w-full overflow-hidden shadow", provider.builtin ? " cursor-not-allowed" : "", showOptionTagAnyhow ? "h-32" : "h-28")}
                    styles={{ body: { height: "100%", padding: "0.5rem 1rem" } }}
                    hoverable
                    onClick={() => {
                      if (provider.builtin) {
                        return;
                      }

                      handleProviderTypeSelect(provider.type);
                    }}
                  >
                    <div className="flex size-full flex-col items-center justify-center gap-3 overflow-hidden p-2">
                      <div className="flex items-center justify-center">
                        <Avatar className="bg-stone-100" icon={<img src={provider.icon} />} shape="square" size={32} />
                      </div>
                      <div className="w-full overflow-hidden text-center">
                        <div className={mergeCls("w-full truncate", { "mb-1": showOptionTagAnyhow })}>
                          <Tooltip title={t(provider.name)} mouseEnterDelay={1}>
                            <Typography.Text type={provider.builtin ? "secondary" : void 0}>{t(provider.name) || "\u00A0"}</Typography.Text>
                          </Tooltip>
                        </div>
                        <Show when={showOptionTagAnyhow}>
                          <div className="origin-top scale-80 whitespace-nowrap" style={{ marginInlineEnd: "-8px" }}>
                            <Show when={provider.builtin}>
                              <Tag>{t("access.props.provider.builtin")}</Tag>
                            </Show>
                            <Show when={showOptionTagForDNS && provider.usages.includes(ACCESS_USAGES.DNS)}>
                              <Tag color="#d93f0b99">{t("access.props.provider.usage.dns")}</Tag>
                            </Show>
                            <Show when={showOptionTagForHosting && provider.usages.includes(ACCESS_USAGES.HOSTING)}>
                              <Tag color="#0052cc99">{t("access.props.provider.usage.hosting")}</Tag>
                            </Show>
                            <Show when={showOptionTagForCA && provider.usages.includes(ACCESS_USAGES.CA)}>
                              <Tag color="#0e8a1699">{t("access.props.provider.usage.ca")}</Tag>
                            </Show>
                            <Show when={showOptionTagForNotification && provider.usages.includes(ACCESS_USAGES.NOTIFICATION)}>
                              <Tag color="#1d76db99">{t("access.props.provider.usage.notification")}</Tag>
                            </Show>
                          </div>
                        </Show>
                      </div>
                    </div>
                  </Card>
                </div>
              );
            })}
          </div>
        </Show>
      </div>
    </div>
  );
};

export default AccessProviderPicker;
