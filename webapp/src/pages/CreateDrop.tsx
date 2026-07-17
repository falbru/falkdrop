import { useState } from "react";
import { useNavigate } from "react-router";
import Button from "../components/ui/Button";
import ItemGroup from "../components/ui/ItemGroup";
import { useAuthenticatedMutation } from "../features/auth/hooks";
import createDrop from "../features/drop/api/createDrop";
import createAndUploadResource from "../features/drop/api/createAndUploadResource";
import ResourceDropZone from "../features/drop/components/ResourceDropZone";
import ResourceItem from "../features/drop/components/ResourceItem";
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
    LocalResource,
    unknown,
    Error
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
    { resource_ids: string[] },
    unknown,
    Error
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
    <div className="flex flex-col gap-4 items-stretch">
      <ResourceDropZone onDrop={handleOnDrop} />

      {uploadedResources.length > 0 && (
        <ItemGroup title="Files">
          {uploadedResources.map((res, index) => (
            <ResourceItem
              key={res.id}
              name={res.name ?? res.id}
              size="64 GB"
              showBorder={index < uploadedResources.length - 1}
            />
          ))}
        </ItemGroup>
      )}

      <div className="flex justify-end">
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
