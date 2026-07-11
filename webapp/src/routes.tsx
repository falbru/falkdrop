import { createBrowserRouter } from "react-router";
import Layout from "./components/layout/Layout";
import HomePage from "./pages/Home";
import LoginPage from "./pages/Login";
import CreateDropPage from "./pages/CreateDrop";
import GetDropPage from "./pages/GetDrop";
import ProtectedRoute from "./components/shared/ProtectedRoute";

const router = createBrowserRouter([
  {
    path: "/login",
    Component: LoginPage,
  },
  {
    element: <Layout showHeader={false} />,
    children: [
      {
        index: true,
        Component: HomePage,
      },
    ],
  },
  {
    element: <Layout />,
    children: [
      {
        path: "/create",
        element: (
          <ProtectedRoute>
            <CreateDropPage />
          </ProtectedRoute>
        ),
      },
      {
        path: "/:dropId",
        Component: GetDropPage,
      },
    ],
  },
]);

export default router;
