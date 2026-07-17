import type { ButtonProps as AriaButtonProps } from "react-aria-components";
import { Button as AriaButton } from "react-aria-components";

type ButtonProps = AriaButtonProps & {
  variant?: "primary" | "secondary" | "ghost";
};

const Button = (props: ButtonProps) => {
  const baseStyles =
    "rounded-lg py-2 px-4 text-sm font-medium transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-border focus:ring-offset-2 focus:ring-offset-background disabled:opacity-50 disabled:cursor-not-allowed";

  let variantStyle = "";
  switch (props.variant) {
    case "primary":
      variantStyle =
        "bg-primary text-primary-text hover:bg-primary-hover active:bg-primary-active";
      break;
    case "secondary":
      variantStyle =
        "bg-secondary text-secondary-text border border-secondary hover:bg-secondary-hover active:bg-secondary-active";
      break;
    case "ghost":
    default:
      variantStyle =
        "bg-transparent text-text hover:bg-secondary active:bg-secondary-active";
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
