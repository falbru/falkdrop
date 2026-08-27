import Keycloak from "keycloak-js";
import type { AuthProvider } from "../types";

const keycloak = new Keycloak({
  url: import.meta.env.VITE_KEYCLOAK_URL as string,
  realm: import.meta.env.VITE_KEYCLOAK_REALM as string,
  clientId: import.meta.env.VITE_KEYCLOAK_CLIENT_ID as string,
});

export function createKeycloakAuthProvider(): AuthProvider {
  return {
    init: async () => {
      return keycloak
        .init({
          onLoad: "check-sso",
          silentCheckSsoRedirectUri:
            typeof window !== "undefined"
              ? window.location.origin + "/silent-callback.html"
              : undefined,
        })
        .then((authenticated) => {
          return authenticated;
        })
        .catch((error: unknown) => {
          console.error("Failed to initialize Keycloak:", error);
          return false;
        });
    },
    login: async (options) => {
      await keycloak.login({ redirectUri: options?.redirectUri });
    },
    logout: async (options) => {
      await keycloak.logout({ redirectUri: options?.redirectUri });
    },
    isAuthenticated() {
      return keycloak.authenticated;
    },
    getToken() {
      return keycloak.token;
    },
  };
}
