import type { InputProps as AriaInputProps } from "react-aria-components";
import { Input as AriaInput } from "react-aria-components";

type InputProps = AriaInputProps & {
  variant?: "default" | "ghost";
};

const Input = (props: InputProps) => {
  const baseStyles =
    "w-full rounded-lg py-2 px-3 text-sm transition-colors duration-200 placeholder:text-neutral-500 focus:outline-none focus:ring-2 focus:ring-neutral-500 focus:ring-offset-0 disabled:opacity-50 disabled:cursor-not-allowed";

  let variantStyle = "";
  switch (props.variant) {
    case "ghost":
      variantStyle = "bg-transparent border-none text-neutral-300";
      break;
    case "default":
    default:
      variantStyle =
        "bg-card border border-neutral-700 text-neutral-100 hover:border-neutral-600";
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
