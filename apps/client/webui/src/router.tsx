import {createBrowserRouter} from "react-router-dom";
import LoginPage from "./pages/login";
import ProtectedGuard from "./guards/protected-guard";
import HostsPage from "./pages/hosts";
import AuthLayout from "./layouts/auth-layout";
import RegisterPage from "./pages/register";
import DefaultLayout from "./layouts/default-layout";
import SingleHostPage from "./pages/host";
import SetupGuard from "./guards/setup-guard";

const router = createBrowserRouter([
  {
    path: "*",
    element: <SetupGuard />,
    children: [
      // Auth routes (login/register)
      {
        element: <AuthLayout />,
        children: [
          {
            path: "login",
            element: <LoginPage />,
          },
          {
            path: "register",
            element: <RegisterPage />,
          },
        ],
      },
      // Protected routes
      {
        element: <ProtectedGuard />,
        children: [
          {
            element: <DefaultLayout />,
            children: [
              {
                path: "hosts",
                element: <HostsPage />,
              },
              {
                path: "hosts/:hostId",
                element: <SingleHostPage />,
              },
            ],
          },
        ],
      },
    ],
  },
]);

export default router;
