import { DropZone, FileTrigger, type DropItem } from "react-aria-components";
import type { LocalResource } from "../types";
import { Button } from "../../../components/ui/button";
import { Upload } from "lucide-react";

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
      className="flex items-center justify-center bg-card p-4 w-full h-[200px] rounded-xl border-3 border-dashed border-border"
    >
      <div className="flex flex-col items-center gap-4">
        <FileTrigger
          onSelect={(files) => {
            if (!files) return;

            const file = files[0];
            if (file) {
              onDrop({
                type: "file",
                getBody: () => Promise.resolve(file),
              });
            }
          }}
        >
          <Button variant="secondary" className="px-3 py-3" size="icon-lg">
            <Upload />
          </Button>
        </FileTrigger>
        <div slot="label" className="text-muted-foreground">
          Drop files or paste text here
        </div>
      </div>
    </DropZone>
  );
};

export default ResourceDropZone;
