import { API_URL } from "../../../config/env";
import type { DropWithDownloadableResources } from "../types";
import {
  toDropWithDownloadableResources,
  type APIDropWithDownloadableResources,
} from "./types";

type GetDropResponse = APIDropWithDownloadableResources;

const fetchDrop = async (
  dropId: string,
): Promise<DropWithDownloadableResources> => {
  return fetch(`${API_URL}/drop/${dropId}`, {
    method: "GET",
  }).then(async (response) => {
    const res = (await response.json()) as GetDropResponse;
    return toDropWithDownloadableResources(res);
  });
};

export default fetchDrop;
