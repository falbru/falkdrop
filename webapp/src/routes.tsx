import { createBrowserRouter } from "react-router";
import Layout from "./components/layout/Layout";
import HomePage from "./pages/Home";
import LoginPage from "./pages/Login";
import CreateDropPage from "./pages/CreateDrop";
import GetDropPage from "./pages/GetDrop";

const router = createBrowserRouter([
  {
    path: "/",
    Component: HomePage,
  },
  {
    path: "/login",
    Component: LoginPage,
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
