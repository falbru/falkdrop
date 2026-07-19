import { useAuth } from "../features/auth/contexts/AuthProvider";
import { useEffect } from "react";
import { useLocation } from "react-router";

export interface LocationState {
  redirectUri?: string;
}

export default function LoginPage() {
  const auth = useAuth();
  const location = useLocation();

  useEffect(() => {
    if (auth) {
      const redirectUri =
        (location.state as LocationState | null)?.redirectUri ??
        window.location.origin;
      void auth.login({ redirectUri });
    }
  }, [auth, location.state]);

  return null;
}
