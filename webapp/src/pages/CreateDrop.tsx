import { useState } from "react";
import Button from "../components/ui/Button";
import useCreateResource from "../features/drop/hooks/useCreateResource";
import useCreateDrop from "../features/drop/hooks/useCreateDrop";
import ResourceDropZone from "../features/drop/components/ResourceDropZone";
import type {
  DropWithDownloadableResources,
  LocalResource,
  Resource,
} from "../features/drop/types";
import { useNavigate } from "react-router";

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

  return (
    <>
      <ResourceDropZone onDrop={handleOnDrop} />
      <ul>
        {uploadedResources.map((res) => (
          <li key={res.id}>{res.id}</li>
        ))}
      </ul>
      <Button onClick={() => handleCreateDrop()}>Create drop</Button>
    </>
  );
};

export default CreateDropPage;
