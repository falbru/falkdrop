import { useMutation, type UseMutationOptions } from "@tanstack/react-query";
import { useAuth } from "../contexts/AuthProvider";

export function useAuthenticatedMutation<TData, TVariables, TContext, TError>(
  options: Omit<
    UseMutationOptions<TData, TError, TVariables, TContext>,
    "mutationFn"
  > & {
    mutationFn: (variables: TVariables, token: string) => Promise<TData>;
  },
) {
  const auth = useAuth();

  return useMutation<TData, TError, TVariables, TContext>({
    ...options,
    mutationFn: (variables) => {
      const token = auth.getToken();
      if (!token) {
        throw new Error("No authentication token available");
      }
      return options.mutationFn(variables, token);
    },
  });
}
