import { useMutation } from "@tanstack/react-query";
import type { LocalResource, Resource } from "../types";
import createResource from "../api/createResource";
import uploadResource from "../api/uploadResource";

const createAndUploadResource = async (
  resource: LocalResource,
): Promise<Resource> => {
  const draft = await createResource({
    resource_type: resource.type,
  });

  const success = uploadResource(draft.uploadURL, await resource.getBody());
  if (!success) throw new Error();

  return {
    id: draft.id,
    type: draft.type,
  };
};

type UseCreateResourceProps = {
  onSuccess?: (resource: Resource) => void;
};

const useCreateResource = (props?: UseCreateResourceProps) => {
  return useMutation({
    mutationFn: createAndUploadResource,
    onSuccess: props?.onSuccess,
  });
};

export default useCreateResource;
