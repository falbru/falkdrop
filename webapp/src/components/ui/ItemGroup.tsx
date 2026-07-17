import type { HTMLAttributes, ReactNode } from "react";
import React from "react";

type ItemGroupProps = HTMLAttributes<HTMLUListElement> & {
  title?: string;
  empty?: ReactNode;
  children: ReactNode;
};

const ItemGroup = ({
  title,
  empty,
  children,
  className = "",
  ...props
}: ItemGroupProps) => {
  const hasChildren = React.Children.count(children) > 0;

  return (
    <div className="flex flex-col gap-2">
      {title && (
        <h2 className="text-sm font-medium text-text-secondary">{title}</h2>
      )}
      <ul
        {...props}
        className={`rounded-xl border border-border bg-card ${className}`}
      >
        {hasChildren ? children : empty}
      </ul>
    </div>
  );
};

export default ItemGroup;
