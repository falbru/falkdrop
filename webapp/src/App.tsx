import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import router from "./routes";
import { RouterProvider } from "react-router/dom";
import { AuthProvider } from "./features/auth/contexts/AuthProvider";
import { createKeycloakAuthProvider } from "./features/auth/providers/keycloak";

const queryClient = new QueryClient();
const keycloakProvider = createKeycloakAuthProvider();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider provider={keycloakProvider}>
        <RouterProvider router={router} />
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
