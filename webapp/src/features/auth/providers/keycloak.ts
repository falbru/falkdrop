import Keycloak from "keycloak-js";
import type { AuthProvider } from "../types";

const keycloak = new Keycloak({
  url: import.meta.env.VITE_KEYCLOAK_URL,
  realm: import.meta.env.VITE_KEYCLOAK_REALM,
  clientId: import.meta.env.VITE_KEYCLOAK_CLIENT_ID,
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
          if (authenticated) {
            console.log("User is authenticated");
          } else {
            console.log("User is not authenticated");
          }
          return authenticated;
        })
        .catch((error) => {
          console.error("Failed to initialize Keycloak:", error);
          return false;
        });
    },
    login: async (options) => {
      keycloak.login({ redirectUri: options?.redirectUri });
    },
    logout: async (options) => {
      keycloak.logout({ redirectUri: options?.redirectUri });
    },
    isAuthenticated() {
      return keycloak.authenticated;
    },
    getToken() {
      return keycloak.token;
    },
  };
}
