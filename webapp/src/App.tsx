import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import CreateDropPage from "./pages/CreateDrop";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <CreateDropPage />
    </QueryClientProvider>
  );
}

export default App;
