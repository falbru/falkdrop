import { useMutation } from "@tanstack/react-query";
import type { DropWithDownloadableResources } from "../types";
import createDrop from "../api/createDrop";

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
