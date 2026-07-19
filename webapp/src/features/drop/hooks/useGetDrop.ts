import { skipToken, useQuery } from "@tanstack/react-query";
import getDrop from "../api/getDrop";

const useDrop = (dropId: string | undefined) => {
  return useQuery({
    queryFn: dropId ? () => getDrop(dropId) : skipToken,
    queryKey: ["drop", dropId],
    enabled: dropId !== undefined,
    refetchOnWindowFocus: false,
    retry: 0,
  });
};

export default useDrop;
