import type { LocalResource, Resource } from "../types";
import createResource from "./createResource";
import uploadResource from "./uploadResource";

const createAndUploadResource = async (
  resource: LocalResource,
  token: string,
): Promise<Resource> => {
  const resourceBody = await resource.getBody();

  const draft = await createResource(
    {
      resource_type: resource.type,
      name: resource.type == "file" ? (resourceBody as File).name : null,
    },
    token,
  );

  const success = await uploadResource(draft.uploadURL, resourceBody);
  if (!success) throw new Error();

  return {
    id: draft.id,
    name: draft.name,
    type: draft.type,
  };
};

export default createAndUploadResource;
