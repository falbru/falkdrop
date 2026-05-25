import { createBrowserRouter } from "react-router";
import HomePage from "./pages/Home";
import CreateDropPage from "./pages/CreateDrop";

const router = createBrowserRouter([
  {
    path: "/",
    Component: HomePage,
  },
  {
    path: "/create",
    Component: CreateDropPage,
  },
]);

export default router;
