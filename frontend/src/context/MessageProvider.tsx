import { message } from "antd";
import type { ReactNode } from "react";
import { MessageContext } from "./MessageContext";

interface Props {
  children: ReactNode;
}

const MessageProvider = ({ children }: Props) => {
  const [messageApi, contextHolder] = message.useMessage();

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      messageApi.success("Copied to clipboard", 1.5);
    } catch (e) {
      messageApi.error(`Failed to copy to clipboard: ${e}`, 1.5);
    }
  };

  return (
    <MessageContext.Provider
      value={{
        messageApi: messageApi,
        copyToClipboard: copyToClipboard,
      }}
    >
      {contextHolder}
      {children}
    </MessageContext.Provider>
  );
};

export default MessageProvider;
