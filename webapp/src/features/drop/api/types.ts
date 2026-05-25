import type {
  DownloadableResource,
  ResourceType,
  UploadableResource,
} from "../types";

type APIResourceBase = {
  id: string;
  type: string;
};

export type APIUploadableResource = APIResourceBase & {
  upload_url: string;
};

export type APIDownloadableResource = APIResourceBase & {
  download_url: string;
};

export const isResourceType = (type: string): type is ResourceType => {
  return type === "file" || type === "text";
};

export const toResourceType = (type: string): ResourceType => {
  if (isResourceType(type)) return type;
  throw new Error(`Invalid resource type from API: "${type}"`);
};

export const toUploadableResource = (
  resource: APIUploadableResource,
): UploadableResource => {
  return {
    id: resource.id,
    type: toResourceType(resource.type),
    uploadURL: resource.upload_url,
  };
};

export const toDownloadableResource = (
  resource: APIDownloadableResource,
): DownloadableResource => {
  return {
    id: resource.id,
    type: toResourceType(resource.type),
    downloadURL: resource.download_url,
  };
};

export type APIDrop = {
  id: string;
  expiration_date: string;
};

export type APIDropWithDownloadableResources = APIDrop & {
  resources: APIDownloadableResource[];
};

export const toDropWithDownloadableResources = (
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
