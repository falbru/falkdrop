import type { InputProps as AriaInputProps } from "react-aria-components";
import { Input as AriaInput } from "react-aria-components";

type InputProps = AriaInputProps & {
  variant?: "default" | "ghost";
};

const Input = (props: InputProps) => {
  const baseStyles =
    "w-full rounded-lg py-2 px-3 text-sm transition-colors duration-200 placeholder:text-text-muted focus:outline-none focus:ring-2 focus:ring-border focus:ring-offset-0 disabled:opacity-50 disabled:cursor-not-allowed";

  let variantStyle = "";
  switch (props.variant) {
    case "ghost":
      variantStyle = "bg-transparent border-none text-text";
      break;
    case "default":
    default:
      variantStyle =
        "bg-card border border-border text-text hover:border-text-muted";
      break;
  }

  return (
    <AriaInput
      {...props}
      className={`${variantStyle} ${baseStyles} ${props.className || ""}`}
    />
  );
};

export default Input;
