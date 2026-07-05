import { useAuth } from "../features/auth/contexts/AuthProvider";
import { useEffect } from "react";

export default function LoginPage() {
  const auth = useAuth();

  useEffect(() => {
    if (!auth.isLoading) auth.login({ redirectUri: window.location.origin });
  }, [auth]);

  return null;
}
