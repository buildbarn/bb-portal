import { createLink } from "@tanstack/react-router";
import { Button, type ButtonProps } from "antd";
import * as React from "react";

// Tanstack Router and Antd both have a type prop, so we move the antd type
// prop to `buttonType`.
interface CustomButtonProps extends Omit<ButtonProps, "type"> {
  buttonType?: ButtonProps["type"];
}

// Looks like and can be styles like an AntD button,
// but behaves like a Tanstack link, including type safety and pre-fetching on hover.

// TanStack Router needs the ref to manage focus and accessibility.
const AntdButtonWrapper = React.forwardRef<
  HTMLAnchorElement | HTMLButtonElement,
  CustomButtonProps
>(({ buttonType, ...props }, ref) => {
  return <Button ref={ref} type={buttonType} {...props} />;
});

// Wrap it in createLink to inject all TanStack Router functionality
export const LinkButton = createLink(AntdButtonWrapper);
