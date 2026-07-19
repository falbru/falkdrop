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
    if (response.status >= 400) {
      let errorMessage = "Failed to load drop";

      try {
        const errorData = await response.json();
        errorMessage = errorData.message || errorMessage;
      } catch (e) {
        errorMessage = response.statusText || errorMessage;
      }

      throw new Error(errorMessage);
    }

    const res = (await response.json()) as GetDropResponse;
    return toDropWithDownloadableResources(res);
  });
};

export default fetchDrop;
