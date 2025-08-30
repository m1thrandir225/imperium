import SetupPage from "@/pages/setup";
import useConfigStore from "@/stores/config.store";
import {Loader2} from "lucide-react";
import type React from "react";
import {Navigate, Outlet, useLocation} from "react-router-dom";

const SetupGuard: React.FC = () => {
  const {isConfigured, isLoading, error} = useConfigStore();
  const location = useLocation();

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Loader2 className="w-4 h-4 animate-spin" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-red-600 mb-2">
            Configuration Error
          </h1>
          <p className="text-gray-600">
            Unable to check configuration status. Please try again.
          </p>
        </div>
      </div>
    );
  }

  if (!isConfigured) {
    return <SetupPage />;
  }

  if (isConfigured && location.pathname === "/") {
    return <Navigate to="/hosts" replace />;
  }

  return <Outlet />;
};

export default SetupGuard;
