import type { AuthProvider as AuthProviderType } from "../types";
import { useEffect, useRef, createContext, useContext } from "react";
import type { ReactNode } from "react";

interface AuthContextValue {
  provider: AuthProviderType;
}

const AuthContext = createContext<AuthContextValue | null>(null);

type AuthProviderProps = {
  children: ReactNode;
  provider: AuthProviderType;
};

export function AuthProvider(props: AuthProviderProps) {
  const { children, provider } = props;

  const initialized = useRef(false);

  useEffect(() => {
    if (initialized.current) return; // Ensure provider.init() isn't called twice in StrictMode
    initialized.current = true;
    provider.init();
  }, [provider]);

  return (
    <AuthContext.Provider value={{ provider }}>{children}</AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context.provider;
}
