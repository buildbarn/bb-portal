import { CopyOutlined } from "@ant-design/icons";
import { useBbPortalMessage } from "@/context/MessageContext";

interface Props {
  text: string;
}

const CopyableIcon: React.FC<Props> = ({ text }) => {
  const { copyToClipboard } = useBbPortalMessage();

  return (
    <button
      type="button"
      // Unset the default <button> style
      style={{
        appearance: "none",
        backgroundColor: "transparent",
        border: "none",
        borderRadius: 0,
        font: "inherit",
        lineHeight: 0,
        margin: 0,
        marginRight: "5px",
        padding: 0,
        textAlign: "inherit",
      }}
      onClick={() => copyToClipboard(text)}
    >
      {/** biome-ignore lint/a11y/useValidAnchor: Use an <a> without `href` (which is valid according to HTML specification) for the link styling */}
      <a>
        <CopyOutlined />
      </a>
    </button>
  );
};

export default CopyableIcon;
