import { DropZone, type DropItem } from "react-aria-components";
import type { LocalResource } from "..";

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
    >
      <div slot="label">Drop or paste text or images here</div>
    </DropZone>
  );
};

export default ResourceDropZone;
