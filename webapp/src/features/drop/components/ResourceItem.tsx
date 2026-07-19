import type { ReactNode } from "react";
import { File } from "lucide-react";
import { Item, ItemMedia, ItemContent, ItemTitle } from "@/components/ui/item";

interface ResourceItemProps {
  name: string;
  icon?: ReactNode;
  className?: string;
  children?: ReactNode;
}

const ResourceItem = ({
  name,
  icon,
  className = "",
  children,
  ...props
}: ResourceItemProps) => {
  return (
    <Item className={className} {...props} variant="outline">
      <ItemMedia variant="icon">{icon ?? <File />}</ItemMedia>
      <ItemContent>
        <ItemTitle>{name}</ItemTitle>
      </ItemContent>
      {children && <div className="flex-grow flex justify-end">{children}</div>}
    </Item>
  );
};

export default ResourceItem;
