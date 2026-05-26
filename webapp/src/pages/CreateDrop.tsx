import { useState } from "react";
import { useNavigate } from "react-router";
import Button from "../components/ui/Button";
import Card from "../components/ui/Card";
import useCreateResource from "../features/drop/hooks/useCreateResource";
import useCreateDrop from "../features/drop/hooks/useCreateDrop";
import ResourceDropZone from "../features/drop/components/ResourceDropZone";
import type {
  DropWithDownloadableResources,
  LocalResource,
  Resource,
} from "../features/drop/types";
import { Plus } from "lucide-react";

const CreateDropPage = () => {
  const [uploadedResources, setUploadedResources] = useState<Resource[]>([]);

  const createResourceMutation = useCreateResource({
    onSuccess: (uploadedResource: Resource) => {
      setUploadedResources([...uploadedResources, uploadedResource]);
    },
  });

  const handleOnDrop = (resource: LocalResource) => {
    createResourceMutation.mutate(resource);
  };

  const navigate = useNavigate();
  const createDropMutation = useCreateDrop({
    onSuccess: (drop: DropWithDownloadableResources) => {
      navigate(`/${drop.id}`);
    },
  });

  const handleCreateDrop = () => {
    createDropMutation.mutate({
      resource_ids: uploadedResources.map((res) => res.id),
    });
  };

  const isLoading =
    createResourceMutation.isPending || createDropMutation.isPending;

  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-4">
      <div className="w-[640px] max-w-screen max-w-2xl flex flex-col gap-4 items-center p-8">
        <ResourceDropZone onDrop={handleOnDrop} />

        {uploadedResources.length > 0 && (
          <Card className="w-full">
            <h2 className="text-sm font-medium text-neutral-400 mb-3">
              Uploaded Files
            </h2>
            <ul className="space-y-2">
              {uploadedResources.map((res) => (
                <li
                  key={res.id}
                  className="text-sm text-neutral-300 truncate max-w-full"
                >
                  {res.id}
                </li>
              ))}
            </ul>
          </Card>
        )}

        <div className="flex justify-center">
          <Button
            onClick={handleCreateDrop}
            variant="primary"
            isDisabled={uploadedResources.length === 0 || isLoading}
          >
            {isLoading ? (
              "Creating..."
            ) : (
              <div className="flex gap-1 items-center">
                <Plus size={18} />
                <span>Create Drop</span>
              </div>
            )}
          </Button>
        </div>
      </div>
    </div>
  );
};

export default CreateDropPage;
