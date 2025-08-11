import { useEffect, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useSize } from "ahooks";
import { Avatar, Card, Checkbox, Empty, Flex, Input, type InputRef, Tabs, Tooltip, Typography } from "antd";

import Show from "@/components/Show";
import { DEPLOYMENT_CATEGORIES, type DeploymentProvider, deploymentProvidersMap } from "@/domain/provider";
import { useZustandShallowSelector } from "@/hooks";
import { useAccessesStore } from "@/stores/access";
import { mergeCls } from "@/utils/css";

export interface DeploymentProviderPickerProps {
  className?: string;
  style?: React.CSSProperties;
  autoFocus?: boolean;
  gap?: number | "small" | "middle" | "large";
  placeholder?: string;
  onFilter?: (value: string, option: DeploymentProvider) => boolean;
  onSelect?: (value: string) => void;
}

const DeploymentProviderPicker = ({ className, style, autoFocus, onFilter, placeholder, onSelect, ...props }: DeploymentProviderPickerProps) => {
  const { gap = "middle" } = props;

  const { t } = useTranslation();

  const { accesses, fetchAccesses } = useAccessesStore(useZustandShallowSelector(["accesses", "fetchAccesses"]));
  useEffect(() => {
    fetchAccesses(false);
  }, []);

  const wrapperRef = useRef<HTMLDivElement>(null);
  const wrapperSize = useSize(wrapperRef);

  const [isAvailableOnly, setIsAvailableOnly] = useState(true);

  const [category, setCategory] = useState<string>(DEPLOYMENT_CATEGORIES.ALL);

  const [keyword, setKeyword] = useState<string>();
  const keywordInputRef = useRef<InputRef>(null);
  useEffect(() => {
    if (autoFocus) {
      setTimeout(() => keywordInputRef.current?.focus(), 1);
    }
  }, []);

  const providers = useMemo(() => {
    return Array.from(deploymentProvidersMap.values())
      .filter((provider) => {
        if (onFilter) {
          return onFilter(provider.type, provider);
        }

        return true;
      })
      .filter((provider) => {
        if (isAvailableOnly) {
          return provider.builtin || accesses.some((access) => access.provider === provider.provider);
        }

        return true;
      })
      .filter((provider) => {
        if (category && category !== DEPLOYMENT_CATEGORIES.ALL) {
          return provider.category === category;
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
  }, [onFilter, accesses, isAvailableOnly, category, keyword]);
  const providerCols = useMemo(() => {
    if (!wrapperSize) {
      return 1;
    }

    const cols = Math.floor(wrapperSize.width / 320);
    return Math.min(9, Math.max(1, cols));
  }, [wrapperSize]);

  const handleProviderTypeSelect = (value: string) => {
    onSelect?.(value);
  };

  return (
    <div className={className} style={style} ref={wrapperRef}>
      <Input.Search ref={keywordInputRef} placeholder={placeholder ?? t("common.text.search")} onChange={(e) => setKeyword(e.target.value.trim())} />

      <div className="mt-4">
        <Flex justify="end">
          <Checkbox checked={isAvailableOnly} onClick={() => setIsAvailableOnly(!isAvailableOnly)}>
            {t("provider.text.show_available_hosting_provider_only")}
          </Checkbox>
        </Flex>
      </div>

      <div className="mt-4">
        <Flex>
          <Tabs
            defaultActiveKey={DEPLOYMENT_CATEGORIES.ALL}
            items={[
              DEPLOYMENT_CATEGORIES.ALL,
              DEPLOYMENT_CATEGORIES.CDN,
              DEPLOYMENT_CATEGORIES.STORAGE,
              DEPLOYMENT_CATEGORIES.LOADBALANCE,
              DEPLOYMENT_CATEGORIES.FIREWALL,
              DEPLOYMENT_CATEGORIES.AV,
              DEPLOYMENT_CATEGORIES.ACCELERATOR,
              DEPLOYMENT_CATEGORIES.APIGATEWAY,
              DEPLOYMENT_CATEGORIES.SERVERLESS,
              DEPLOYMENT_CATEGORIES.WEBSITE,
              DEPLOYMENT_CATEGORIES.SSL,
              DEPLOYMENT_CATEGORIES.NAS,
              DEPLOYMENT_CATEGORIES.OTHER,
            ].map((key) => ({
              key: key,
              label: t(`provider.category.${key}`),
            }))}
            size="small"
            tabBarStyle={{ marginLeft: "-1rem" }}
            tabPosition="left"
            onChange={(key) => setCategory(key)}
          />

          <div className="flex-1">
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
        </Flex>
      </div>
    </div>
  );
};

export default DeploymentProviderPicker;
