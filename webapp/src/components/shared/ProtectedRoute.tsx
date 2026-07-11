import { useEffect, type ReactNode } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../../features/auth/contexts/AuthProvider";
import type { LocationState } from "../../pages/Login";

type ProtectedRouteProps = {
  children: ReactNode;
};

const ProtectedRoute = (props: ProtectedRouteProps) => {
  const auth = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (auth && !auth.isAuthenticated()) {
      const currentUrl = window.location.href;
      navigate("/login", {
        state: { redirectUri: currentUrl } as LocationState,
      });
    }
  }, [auth, navigate]);

  return auth && auth.isAuthenticated() ? props.children : null;
};

export default ProtectedRoute;
