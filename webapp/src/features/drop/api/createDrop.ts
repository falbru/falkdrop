import { API_URL } from "../../../config/env";
import type { DropWithDownloadableResources } from "../types";
import {
  toDropWithDownloadableResources,
  type APIDropWithDownloadableResources,
} from "./types";

interface CreateDropRequest {
  resource_ids: string[];
}

type CreateDropResponse = APIDropWithDownloadableResources;

const createDrop = async (request: CreateDropRequest, token: string) => {
  const drop: DropWithDownloadableResources = await fetch(`${API_URL}/drop`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(request),
  }).then(async (response) => {
    const drop = (await response.json()) as CreateDropResponse;
    return toDropWithDownloadableResources(drop);
  });

  return drop;
};

export default createDrop;
