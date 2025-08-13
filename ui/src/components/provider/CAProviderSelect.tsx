import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { Avatar, Select, type SelectProps, Typography, theme } from "antd";

import { type CAProvider, caProvidersMap } from "@/domain/provider";

export interface CAProviderSelectProps
  extends Omit<SelectProps, "filterOption" | "filterSort" | "labelRender" | "options" | "optionFilterProp" | "optionLabelProp" | "optionRender"> {
  onFilter?: (value: string, option: CAProvider) => boolean;
}

const CAProviderSelect = ({ onFilter, ...props }: CAProviderSelectProps) => {
  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();

  const options = useMemo<Array<{ key: string; value: string; label: string; data: CAProvider }>>(() => {
    const temp = Array.from(caProvidersMap.values())
      .filter((provider) => {
        if (onFilter) {
          return onFilter(provider.type, provider);
        }

        return true;
      })
      .map((provider) => ({
        key: provider.type as string,
        value: provider.type as string,
        label: t(provider.name),
        data: provider,
      }));

    temp.unshift({
      key: "",
      value: "",
      label: t("provider.text.default_ca_provider"),
      data: {} as CAProvider,
    });

    return temp;
  }, [onFilter]);

  const renderOption = (key: string) => {
    if (key === "") {
      return (
        <div className="flex items-center gap-2 truncate overflow-hidden">
          <Typography.Text className="italic" ellipsis italic>
            {t("provider.text.default_ca_provider")}
          </Typography.Text>
        </div>
      );
    }

    const provider = caProvidersMap.get(key);
    return (
      <div className="flex items-center gap-2 truncate overflow-hidden">
        <Avatar shape="square" src={provider?.icon} size="small" />
        <Typography.Text ellipsis>{t(provider?.name ?? "")}</Typography.Text>
      </div>
    );
  };

  return (
    <Select
      {...props}
      filterOption={(inputValue, option) => {
        if (!option) return false;

        const value = inputValue.toLowerCase();
        return option.value.toLowerCase().includes(value) || option.label.toLowerCase().includes(value);
      }}
      labelRender={({ value }) => {
        if (value != null) {
          return renderOption(value as string);
        }

        return <span style={{ color: themeToken.colorTextPlaceholder }}>{props.placeholder}</span>;
      }}
      options={options}
      optionFilterProp={void 0}
      optionLabelProp={void 0}
      optionRender={(option) => renderOption(option.data.value)}
    />
  );
};

export default CAProviderSelect;
