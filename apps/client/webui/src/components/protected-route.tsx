import useAuthStore from "@/stores/auth.store";
import type React from "react";
import {Navigate, Outlet} from "react-router-dom";

const ProtectedRoute: React.FC = () => {
  const {checkAuth, isHydrated} = useAuthStore();

  if (!isHydrated) {
    return <div>Loading ... </div>;
  }

  return checkAuth() ? <Outlet /> : <Navigate to={"/login"} />;
};

export default ProtectedRoute;
