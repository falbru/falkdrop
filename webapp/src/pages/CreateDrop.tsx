import { useState } from "react";
import { useNavigate } from "react-router";
import { Button } from "../components/ui/button";
import { ItemGroup } from "../components/ui/item";
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
import { toast } from "sonner";

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
    onError: (error: Error) => {
      toast.error(error.message || "Failed to upload resource");
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
    onError: (error: Error) => {
      toast.error(error.message || "Failed to create drop");
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
          {uploadedResources.map((res) => (
            <ResourceItem key={res.id} name={res.name ?? res.id} />
          ))}
        </ItemGroup>
      )}

      <div className="flex justify-end">
        <Button
          onClick={handleCreateDrop}
          variant="default"
          isDisabled={uploadedResources.length === 0 || isLoading}
        >
          {isLoading ? (
            "Creating..."
          ) : (
            <>
              <Plus size={18} />
              <span>Create Drop</span>
            </>
          )}
        </Button>
      </div>
    </div>
  );
};

export default CreateDropPage;
