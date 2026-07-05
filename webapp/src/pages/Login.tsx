import { useAuth } from "../features/auth/contexts/AuthProvider";
import { useEffect } from "react";
import { useNavigate } from "react-router";

export default function LoginPage() {
  const navigate = useNavigate();
  const auth = useAuth();

  useEffect(() => {
    if (auth.isAuthenticated()) {
      navigate("/");
      return;
    }

    auth.login({ redirectUri: window.location.origin });
  }, [navigate, auth]);

  return null;
}
