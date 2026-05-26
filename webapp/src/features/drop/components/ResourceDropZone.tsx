import { DropZone, type DropItem } from "react-aria-components";
import type { LocalResource } from "../types";

type ResourceDropZoneProps = {
  onDrop: (resource: LocalResource) => void;
};

const ResourceDropZone = (props: ResourceDropZoneProps) => {
  const { onDrop } = props;

  const handleItemOnDrop = (item: DropItem) => {
    switch (item.kind) {
      case "text":
        onDrop({
          type: "text",
          getBody: () => item.getText("text/plain"),
        });
        break;
      case "file":
        onDrop({
          type: "file",
          getBody: () => item.getFile(),
        });
        break;
      default:
        throw new Error(`Unsupported drop type: ${item.kind}`);
    }
  };

  return (
    <DropZone
      getDropOperation={() => "copy"}
      onDrop={(event) => {
        event.items.forEach((item) => handleItemOnDrop(item));
      }}
      className="flex items-center justify-center bg-card p-4 w-full h-[250px] rounded-xl border-4 border-dashed border-border"
    >
      <div slot="label" className="text-neutral-700">
        Drop files or paste text here
      </div>
    </DropZone>
  );
};

export default ResourceDropZone;
