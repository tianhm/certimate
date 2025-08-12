import { IconBulb } from "@tabler/icons-react";
import { Alert, Flex, Typography, theme } from "antd";

export interface TipsProps {
  className?: string;
  style?: React.CSSProperties;
  message: React.ReactNode;
}

const Tips = ({ className, style, message }: TipsProps) => {
  const { token: themeToken } = theme.useToken();

  return (
    <Alert
      className={className}
      style={style}
      message={
        <Flex gap="small">
          <div>
            <IconBulb size={18} color={themeToken.colorInfo} />
          </div>
          <Typography.Text>{message}</Typography.Text>
        </Flex>
      }
      type="info"
    />
  );
};

export default Tips;
