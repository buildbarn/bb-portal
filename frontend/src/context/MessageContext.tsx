import type { MessageInstance } from "antd/lib/message/interface";
import { createContext, useContext } from "react";

interface MessageContextState {
  messageApi: MessageInstance;
  copyToClipboard: (text: string) => void;
}

// biome-ignore lint/style/noNonNullAssertion: We want to throw an error if the context is used without provider, instead of failing silently.
export const MessageContext = createContext<MessageContextState>(null!);

export const useBbPortalMessage = () => useContext(MessageContext);
