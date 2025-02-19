import themeStyles from "@/theme/theme.module.css";
import { Tag } from "antd";

interface Props {
  color?: string;
  children: React.ReactNode;
}

const SizeClassOutcomeTag: React.FC<Props> = ({
  color = "default",
  children,
}) => {
  return (
    <Tag color={color} className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>{children}</span>
    </Tag>
  );
};

export default SizeClassOutcomeTag;
