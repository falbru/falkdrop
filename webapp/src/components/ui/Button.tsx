import type { ButtonProps as AriaButtonProps } from "react-aria-components";
import { Button as AriaButton } from "react-aria-components";

type ButtonProps = AriaButtonProps & {
  variant?: "primary" | "secondary" | "ghost";
};

const Button = (props: ButtonProps) => {
  const baseStyles =
    "rounded-lg py-2 px-4 text-sm font-medium transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-neutral-500 focus:ring-offset-2 focus:ring-offset-neutral-950 disabled:opacity-50 disabled:cursor-not-allowed";

  let variantStyle = "";
  switch (props.variant) {
    case "primary":
      variantStyle =
        "bg-neutral-100 text-neutral-900 hover:bg-neutral-200 active:bg-neutral-300";
      break;
    case "secondary":
      variantStyle =
        "bg-neutral-800 text-neutral-300 border border-neutral-700 hover:bg-neutral-700 active:bg-neutral-600";
      break;
    case "ghost":
    default:
      variantStyle =
        "bg-transparent text-neutral-300 hover:bg-neutral-800 active:bg-neutral-700";
      break;
  }

  return (
    <AriaButton
      {...props}
      className={`${variantStyle} ${baseStyles} ${props.className || ""}`}
    />
  );
};

export default Button;
