import type React from "react";

const BuildbarnIconSvg: React.FC<React.SVGProps<SVGSVGElement>> = (props) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    version="1.1"
    viewBox="-0.5 -0.5 421 421"
    {...props}
  >
    <title>{props.name}</title>
    <defs />
    <g>
      <rect
        x="180"
        y="120"
        width="60"
        height="60"
        fill="none"
        stroke={props.stroke || "rgb(0, 0, 0)"}
        stroke-width="20"
        pointer-events="all"
      />
      <path
        stroke-linecap="round"
        d="M 10 250 L 50 90 L 210 11 L 370 90 L 410 250"
        fill="none"
        stroke={props.stroke || "rgb(0, 0, 0)"}
        stroke-width="20"
        stroke-miterlimit="10"
        pointer-events="stroke"
      />
      <path
        stroke-linecap="round"
        d="M 70 410 L 70 210"
        fill="none"
        stroke={props.stroke || "rgb(0, 0, 0)"}
        stroke-width="20"
        stroke-miterlimit="10"
        pointer-events="stroke"
      />
      <path
        stroke-linecap="round"
        d="M 350 410 L 350 210"
        fill="none"
        stroke={props.stroke || "rgb(0, 0, 0)"}
        stroke-width="20"
        stroke-miterlimit="10"
        pointer-events="stroke"
      />
      <path
        stroke-linecap="round"
        d="M 10 410 L 410 410"
        fill="none"
        stroke={props.stroke || "rgb(0, 0, 0)"}
        stroke-width="20"
        stroke-miterlimit="10"
        pointer-events="stroke"
      />
      <rect
        x="140"
        y="270"
        width="140"
        height="140"
        fill="none"
        stroke={props.stroke || "rgb(0, 0, 0)"}
        stroke-width="20"
        pointer-events="all"
      />
      <path
        stroke-linecap="round"
        d="M 140 410 L 280 270"
        fill="none"
        stroke={props.stroke || "rgb(0, 0, 0)"}
        stroke-width="20"
        stroke-miterlimit="10"
        pointer-events="stroke"
      />
      <path
        stroke-linecap="round"
        d="M 140 270 L 280 410"
        fill="none"
        stroke={props.stroke || "rgb(0, 0, 0)"}
        stroke-width="20"
        stroke-miterlimit="10"
        pointer-events="stroke"
      />
    </g>
  </svg>
);

export default BuildbarnIconSvg;
