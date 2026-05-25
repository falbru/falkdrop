import { API_URL } from "../../../config/env";
import type { UploadableResource } from "../types";
import { toUploadableResource, type APIUploadableResource } from "./types";

type CreateResourceRequest = {
  resource_type: string;
};
type CreateResourceResponse = APIUploadableResource;

const createResource = async (
  request: CreateResourceRequest,
): Promise<UploadableResource> => {
  return fetch(`${API_URL}/resource`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(request),
  }).then(async (response) => {
    const res = (await response.json()) as CreateResourceResponse;
    return toUploadableResource(res);
  });
};

export default createResource;
