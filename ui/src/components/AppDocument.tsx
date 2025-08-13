import { useTranslation } from "react-i18next";
import { IconBook } from "@tabler/icons-react";
import { Typography } from "antd";

import { APP_DOCUMENT_URL } from "@/domain/app";

export interface AppDocumentLinkButtonProps {
  className?: string;
  style?: React.CSSProperties;
  showIcon?: boolean;
}

const AppDocumentLinkButton = ({ className, style, showIcon = true }: AppDocumentLinkButtonProps) => {
  const { i18n, t } = useTranslation();

  const handleDocumentClick = () => {
    if (i18n.language.startsWith("en")) {
      window.open(APP_DOCUMENT_URL + "/en/", "_blank");
    } else {
      window.open(APP_DOCUMENT_URL, "_blank");
    }
  };

  return (
    <Typography.Link className={className} style={style} type="secondary" onClick={handleDocumentClick}>
      <div className="flex items-center justify-center space-x-1">
        {showIcon ? <IconBook size="1em" /> : <></>}
        <span>{t("common.menu.document")}</span>
      </div>
    </Typography.Link>
  );
};

export default {
  LinkButton: AppDocumentLinkButton,
};
