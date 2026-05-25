import { useMutation } from "@tanstack/react-query";
import type { DropWithDownloadableResources } from "../types";
import {
  toDownloadableResource,
  type APIDownloadableResource,
} from "../api/types";
import { API_URL } from "../../../config/env";

type CreateDropRequest = {
  resource_ids: string[];
};

type APIDrop = {
  id: string;
  expiration_date: string;
};

type APIDropWithDownloadableResources = APIDrop & {
  resources: APIDownloadableResource[];
};
type CreateDropResponse = APIDropWithDownloadableResources;

const toDropWithDownloadableResources = (
  drop: APIDropWithDownloadableResources,
): DropWithDownloadableResources => {
  return {
    id: drop["id"],
    expirationDate: new Date(drop["expiration_date"]),
    resources: drop["resources"].map((resource) =>
      toDownloadableResource(resource),
    ),
  };
};

const createDrop = async (request: CreateDropRequest) => {
  const drop: DropWithDownloadableResources = await fetch(`${API_URL}/drop`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(request),
  }).then(async (response) => {
    const drop = (await response.json()) as CreateDropResponse;
    return toDropWithDownloadableResources(drop);
  });

  return drop;
};

type CreateDropProps = {
  onSuccess?: (drop: DropWithDownloadableResources) => void;
};

const useCreateDrop = (props?: CreateDropProps) => {
  return useMutation({
    mutationFn: createDrop,
    onSuccess: props?.onSuccess,
  });
};

export default useCreateDrop;
