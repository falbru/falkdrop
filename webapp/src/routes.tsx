import { createBrowserRouter } from "react-router";
import Layout from "./components/layout/Layout";
import HomePage from "./pages/Home";
import CreateDropPage from "./pages/CreateDrop";
import GetDropPage from "./pages/GetDrop";

const router = createBrowserRouter([
  {
    path: "/",
    Component: HomePage,
  },
  {
    element: <Layout />,
    children: [
      {
        path: "/create",
        Component: CreateDropPage,
      },
      {
        path: "/:dropId",
        Component: GetDropPage,
      },
    ],
  },
]);

export default router;
