import { useEffect, type ReactNode } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../../features/auth/contexts/AuthProvider";

interface ProtectedRouteProps {
  children: ReactNode;
}

const ProtectedRoute = (props: ProtectedRouteProps) => {
  const auth = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (auth && !auth.isAuthenticated()) {
      const currentUrl = window.location.href;

      void navigate("/login", {
        state: { redirectUri: currentUrl },
      });
    }
  }, [auth, navigate]);

  return auth?.isAuthenticated() ? props.children : null;
};

export default ProtectedRoute;
