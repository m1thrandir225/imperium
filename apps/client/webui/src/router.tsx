import {createBrowserRouter} from "react-router-dom";
import LoginPage from "./pages/login";
import ProtectedRoute from "./components/protected-route";
import HostsPage from "./pages/hosts";

const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        path: "/hosts",
        element: <HostsPage />,
      },
      {
        path: "/",
        element: <div>Home</div>,
      },
    ],
  },
]);

export default router;
