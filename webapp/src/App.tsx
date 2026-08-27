import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import router from "./routes";
import { RouterProvider } from "react-router/dom";
import { AuthProvider } from "./features/auth/contexts/AuthProvider";
import { createKeycloakAuthProvider } from "./features/auth/providers/keycloak";
import { ThemeProvider } from "./components/shared/ThemeProvider";
import { Toaster } from "./components/ui/sonner";

const queryClient = new QueryClient();
const keycloakProvider = createKeycloakAuthProvider();

function App() {
  return (
    <ThemeProvider storageKey="vite-ui-theme">
      <QueryClientProvider client={queryClient}>
        <AuthProvider provider={keycloakProvider}>
          <Toaster />
          <RouterProvider router={router} />
        </AuthProvider>
      </QueryClientProvider>
    </ThemeProvider>
  );
}

export default App;
