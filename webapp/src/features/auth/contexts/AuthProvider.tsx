import type { AuthProvider as AuthProviderType } from "../types";
import { useEffect, useRef, useState, createContext, use } from "react";
import type { ReactNode } from "react";

interface AuthContextValue {
  provider: AuthProviderType;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextValue | null>(null);

interface AuthProviderProps {
  children: ReactNode;
  provider: AuthProviderType;
}

export function AuthProvider(props: AuthProviderProps) {
  const { children, provider } = props;

  const [isLoading, setIsLoading] = useState(true);
  const initialized = useRef(false);

  useEffect(() => {
    if (initialized.current) return; // Ensure provider.init() isn't called twice in StrictMode
    initialized.current = true;
    void provider.init().finally(() => {
      setIsLoading(false);
    });
  }, [provider]);

  return <AuthContext value={{ provider, isLoading }}>{children}</AuthContext>;
}

export function useAuth() {
  const context = use(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return !context.isLoading ? context.provider : null;
}
