import {createBrowserRouter} from "react-router-dom";
import LoginPage from "./pages/login";
import ProtectedRoute from "./components/protected-route";
import HostsPage from "./pages/hosts";
import AuthLayout from "./layouts/auth-layout";
import RegisterPage from "./pages/register";
import DefaultLayout from "./layouts/default-layout";
import SingleHostPage from "./pages/host";

const router = createBrowserRouter([
  {
    element: <AuthLayout />,
    children: [
      {
        path: "/login",
        element: <LoginPage />,
      },
      {
        path: "/register",
        element: <RegisterPage />,
      },
    ],
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        element: <DefaultLayout />,
        children: [
          {
            path: "/hosts",
            element: <HostsPage />,
          },
          {
            path: "/hosts/:hostId",
            element: <SingleHostPage />,
          },
        ],
      },
    ],
  },
]);

export default router;
