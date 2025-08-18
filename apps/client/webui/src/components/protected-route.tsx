import type React from "react";
import {useMemo} from "react";
import {Navigate, Outlet} from "react-router-dom";

const ProtectedRoute: React.FC = () => {
  const isAuthenticated = useMemo(() => {
    return true;
  }, []);
  //TODO: check if user is authenticated

  return isAuthenticated ? <Outlet /> : <Navigate to={"/login"} />;
};

export default ProtectedRoute;
