export type ResourceType = "text" | "file";

export type LocalResourceBody = string | File;

export interface LocalResource {
  type: ResourceType;
  getBody: () => Promise<LocalResourceBody>;
}

export interface Resource {
  id: string;
  type: ResourceType;
  name: string | null;
}

export type DownloadableResource = Resource & {
  downloadURL: string;
};

export type UploadableResource = Resource & {
  uploadURL: string;
};

export interface Drop {
  id: string;
  expirationDate: Date;
}

export type DropWithDownloadableResources = Drop & {
  resources: DownloadableResource[];
};
