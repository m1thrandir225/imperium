import {createBrowserRouter} from "react-router-dom";
import LoginPage from "./pages/login";
import ProtectedGuard from "./guards/protected-guard";
import HostsPage from "./pages/hosts";
import AuthLayout from "./layouts/auth-layout";
import RegisterPage from "./pages/register";
import DefaultLayout from "./layouts/default-layout";
import SingleHostPage from "./pages/host";
import SetupGuard from "./guards/setup-guard";
import SessionPage from "./pages/session";

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
            handle: {
              title: () => "Login | Imperium",
            },
          },
          {
            path: "register",
            element: <RegisterPage />,
            handle: {
              title: () => "Register | Imperium",
            },
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
                handle: {
                  title: () => "Hosts | Imperium",
                },
              },
              {
                path: "hosts/:hostId",
                element: <SingleHostPage />,
                handle: {
                  title: () => "Host Details | Imperium",
                },
              },
              {
                path: "sessions/:sessionId",
                element: <SessionPage />,
                handle: {
                  title: () => "Session | Imperium",
                },
              },
            ],
          },
        ],
      },
    ],
  },
]);

export default router;
