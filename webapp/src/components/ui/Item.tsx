import type { HTMLAttributes } from "react";

type ItemProps = HTMLAttributes<HTMLLIElement> & {
  showBorder?: boolean;
};

const Item = ({
  showBorder = true,
  className = "",
  children,
  ...props
}: ItemProps) => {
  return (
    <li
      {...props}
      className={`text-sm text-text truncate max-w-full py-3 px-4 flex gap-3 items-center ${showBorder ? "border-b border-border" : ""} ${className}`}
    >
      {children}
    </li>
  );
};

export default Item;
