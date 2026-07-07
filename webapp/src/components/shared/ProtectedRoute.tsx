import { useEffect, type ReactNode } from "react";
import { useLocation, useNavigate } from "react-router";
import { useAuth } from "../../features/auth/contexts/AuthProvider";
import type { LocationState } from "../../pages/Login";

type ProtectedRouteProps = {
  children: ReactNode;
};

const ProtectedRoute = (props: ProtectedRouteProps) => {
  const { isLoading, isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    if (!isLoading && !isAuthenticated()) {
      const currentUrl = window.location.href;
      navigate("/login", {
        state: { redirectUri: currentUrl } as LocationState,
      });
    }
  }, [isLoading, isAuthenticated, navigate]);

  return !isLoading && isAuthenticated() ? props.children : null;
};

export default ProtectedRoute;
