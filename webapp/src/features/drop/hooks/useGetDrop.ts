import { useQuery } from "@tanstack/react-query";
import fetchDrop from "../api/fetchDrop";

const useDrop = (dropId?: string) => {
  return useQuery({
    queryFn: () => fetchDrop(dropId!),
    queryKey: ["drop", dropId],
    enabled: dropId !== undefined,
    refetchOnWindowFocus: false,
  });
};

export default useDrop;
