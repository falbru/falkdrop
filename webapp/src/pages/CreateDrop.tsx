import { useState } from "react";
import { useNavigate } from "react-router";
import { Button } from "../components/ui/button";
import { ItemGroup } from "../components/ui/item";
import { useAuthenticatedMutation } from "../features/auth/hooks";
import createDrop, {
  type CreateDropRequest,
} from "../features/drop/api/createDrop";
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
import ExpiryDurationSelect, {
  type Duration,
} from "@/features/drop/components/ExpiryDurationSelect";
import { Label } from "@/components/ui/label";

const CreateDropPage = () => {
  const [uploadedResources, setUploadedResources] = useState<Resource[]>([]);
  const [expiryDuration, setExpiryDuration] = useState<Duration>("PT10M");

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

  const handleSetExpiryDuration = (duration: Duration) => {
    setExpiryDuration(duration);
  };

  const handleOnDrop = (resource: LocalResource) => {
    createResourceMutation.mutate(resource);
  };

  const navigate = useNavigate();
  const createDropMutation = useAuthenticatedMutation<
    DropWithDownloadableResources,
    CreateDropRequest,
    unknown,
    Error
  >({
    mutationFn: createDrop,
    onSuccess: (drop: DropWithDownloadableResources) => {
      void navigate(`/${drop.id}`);
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to create drop");
    },
  });

  const handleCreateDrop = () => {
    createDropMutation.mutate({
      resource_ids: uploadedResources.map((res) => res.id),
      expiry_duration: expiryDuration,
    });
  };

  const isLoading =
    createResourceMutation.isPending || createDropMutation.isPending;

  return (
    <div className="flex flex-col gap-6 items-stretch">
      <ResourceDropZone onDrop={handleOnDrop} />

      {uploadedResources.length > 0 && (
        <ItemGroup title="Files">
          {uploadedResources.map((res) => (
            <ResourceItem key={res.id} name={res.name ?? res.id} />
          ))}
        </ItemGroup>
      )}

      <div className="flex flex-col gap-2">
        <Label>Expire In</Label>
        <ExpiryDurationSelect onSelectionChange={handleSetExpiryDuration} />
      </div>

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
