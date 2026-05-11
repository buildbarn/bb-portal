import { createLink } from "@tanstack/react-router";
import type { ButtonProps } from "antd";
import { Button } from "antd";
import * as React from "react";

// Looks like and can be styles like an AntD button,
// but behaves like a Tanstack link, including type safety and pre-fetching on hover.

// TanStack Router needs the ref to manage focus and accessibility.
const AntdButtonWrapper = React.forwardRef<
  HTMLAnchorElement | HTMLButtonElement,
  ButtonProps
>((props, ref) => {
  return <Button ref={ref} {...props} />;
});

// Wrap it in createLink to inject all TanStack Router functionality
export const LinkButton = createLink(AntdButtonWrapper);
