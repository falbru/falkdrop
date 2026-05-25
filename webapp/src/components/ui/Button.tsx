import type { ButtonProps } from "react-aria-components";
import { Button as AriaButton } from "react-aria-components";

const Button = (props: ButtonProps) => {
  return <AriaButton {...props} className="bg-gray-500 rounded-full" />;
};

export default Button;
