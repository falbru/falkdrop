import type { HTMLAttributes, ReactNode } from "react";
import { File } from "lucide-react";
import Item from "../../../components/ui/Item";

type ResourceItemProps = HTMLAttributes<HTMLLIElement> & {
  name: string;
  size?: string;
  icon?: ReactNode;
  showBorder?: boolean;
};

const ResourceItem = ({
  name,
  size,
  icon,
  showBorder = true,
  className = "",
  children,
  ...props
}: ResourceItemProps) => {
  return (
    <Item showBorder={showBorder} className={className} {...props}>
      <div className="p-2 bg-border rounded-xl">
        {icon ?? <File className="text-white/50" />}
      </div>
      <div className="flex flex-col min-w-0">
        <span className="font-bold truncate">{name}</span>
        {size && <span className="text-xs text-text-muted">{size}</span>}
      </div>
      <div className="flex-grow flex justify-end">{children}</div>
    </Item>
  );
};

export default ResourceItem;
