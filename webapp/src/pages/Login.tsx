import { useAuth } from "../features/auth/contexts/AuthProvider";
import { useEffect } from "react";
import { useLocation } from "react-router";

export type LocationState = {
  redirectUri?: string;
};

export default function LoginPage() {
  const auth = useAuth();
  const location = useLocation();

  useEffect(() => {
    if (!auth.isLoading) {
      const redirectUri =
        (location.state as LocationState)?.redirectUri ??
        window.location.origin;
      auth.login({ redirectUri });
    }
  }, [auth, location.state]);

  return null;
}
