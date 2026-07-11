import {
  useQuery,
  type UseQueryOptions,
  type QueryKey,
} from "@tanstack/react-query";
import { useAuth } from "../contexts/AuthProvider";

export function useAuthenticatedQuery<TData, TError>(
  options: Omit<UseQueryOptions<TData, TError>, "queryKey" | "queryFn"> & {
    queryKey: QueryKey;
    queryFn: (token: string) => Promise<TData>;
  },
) {
  const auth = useAuth();

  return useQuery<TData, TError>({
    ...options,
    queryFn: () => {
      const token = auth?.getToken();
      if (!token) {
        throw new Error("No authentication token available");
      }
      return options.queryFn(token);
    },
  });
}
