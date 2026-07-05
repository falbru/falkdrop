import { useState } from "react";
import { useNavigate } from "react-router";
import Button from "../components/ui/Button";
import Card from "../components/ui/Card";
import { useAuthenticatedMutation } from "../features/auth/hooks";
import createDrop from "../features/drop/api/createDrop";
import createAndUploadResource from "../features/drop/api/createAndUploadResource";
import ResourceDropZone from "../features/drop/components/ResourceDropZone";
import type {
  DropWithDownloadableResources,
  LocalResource,
  Resource,
} from "../features/drop/types";
import { Plus } from "lucide-react";

const CreateDropPage = () => {
  const [uploadedResources, setUploadedResources] = useState<Resource[]>([]);

  const createResourceMutation = useAuthenticatedMutation<
    Resource,
    LocalResource
  >({
    mutationFn: createAndUploadResource,
    onSuccess: (uploadedResource: Resource) => {
      setUploadedResources([...uploadedResources, uploadedResource]);
    },
  });

  const handleOnDrop = (resource: LocalResource) => {
    createResourceMutation.mutate(resource);
  };

  const navigate = useNavigate();
  const createDropMutation = useAuthenticatedMutation<
    DropWithDownloadableResources,
    { resource_ids: string[] }
  >({
    mutationFn: createDrop,
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
    <div className="flex flex-col gap-4 items-center">
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
                {res.name ?? res.id}
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
  );
};

export default CreateDropPage;
