import { createBrowserRouter } from "react-router";
import HomePage from "./pages/Home";
import CreateDropPage from "./pages/CreateDrop";
import GetDropPage from "./pages/GetDrop";

const router = createBrowserRouter([
  {
    path: "/",
    Component: HomePage,
  },
  {
    path: "/create",
    Component: CreateDropPage,
  },
  {
    path: "/:dropId",
    Component: GetDropPage,
  },
]);

export default router;
