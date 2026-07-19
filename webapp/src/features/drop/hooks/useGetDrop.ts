import { useQuery } from "@tanstack/react-query";
import fetchDrop from "../api/fetchDrop";

const useDrop = (dropId: string | undefined) => {
  return useQuery({
    queryFn: () => fetchDrop(dropId!),
    queryKey: ["drop", dropId],
    enabled: dropId !== undefined,
    refetchOnWindowFocus: false,
    retry: 0,
  });
};

export default useDrop;
