import type { HTMLAttributes } from "react";

type CardProps = HTMLAttributes<HTMLDivElement> & {
  variant?: "default" | "glass";
};

const Card = ({
  variant = "default",
  children,
  className = "",
  ...props
}: CardProps) => {
  const baseStyles =
    "rounded-xl p-6 transition-colors duration-200 border border-border";

  let variantStyle = "";
  switch (variant) {
    case "glass":
      variantStyle = "bg-card/50 backdrop-blur-sm";
      break;
    case "default":
    default:
      variantStyle = "bg-card";
      break;
  }

  return (
    <div {...props} className={`${variantStyle} ${baseStyles} ${className}`}>
      {children}
    </div>
  );
};

export default Card;
