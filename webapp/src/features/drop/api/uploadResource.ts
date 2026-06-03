import type { LocalResourceBody } from "../types";

const uploadResource = async (
  uploadURL: string,
  body: LocalResourceBody,
): Promise<boolean> => {
  return fetch(uploadURL, {
    method: "PUT",
    headers: {
      "Content-Type": "application/octet-stream",
    },
    body,
  }).then((response) => response.ok);
};

export default uploadResource;
