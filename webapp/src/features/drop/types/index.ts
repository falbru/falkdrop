export type ResourceType = "text" | "file";

export type LocalResourceBody = string | File;

export type LocalResource = {
  type: ResourceType;
  getBody: () => Promise<LocalResourceBody>;
};

export type Resource = {
  id: string;
  type: ResourceType;
};

export type DownloadableResource = Resource & {
  downloadURL: string;
};

export type UploadableResource = Resource & {
  uploadURL: string;
};

export type Drop = {
  id: string;
  expirationDate: Date;
};

export type DropWithDownloadableResources = Drop & {
  resources: DownloadableResource[];
};
