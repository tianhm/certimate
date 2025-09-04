import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { Avatar, Select, Typography, theme } from "antd";

import { type DeploymentProvider, deploymentProvidersMap } from "@/domain/provider";

import { type SharedSelectProps } from "./_shared";

export interface DeploymentProviderSelectProps extends SharedSelectProps<DeploymentProvider> {}

const DeploymentProviderSelect = ({ onFilter, ...props }: DeploymentProviderSelectProps) => {
  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();

  const options = useMemo<Array<{ key: string; value: string; label: string; data: DeploymentProvider }>>(() => {
    return Array.from(deploymentProvidersMap.values())
      .filter((provider) => {
        if (onFilter) {
          return onFilter(provider.type, provider);
        }

        return true;
      })
      .map((provider) => ({
        key: provider.type,
        value: provider.type,
        label: t(provider.name),
        data: provider,
      }));
  }, [onFilter]);

  const renderOption = (key: string) => {
    const provider = deploymentProvidersMap.get(key);
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

export default DeploymentProviderSelect;
