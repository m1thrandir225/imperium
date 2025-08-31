import Logo from "@/components/logo";
import {Toaster} from "@/components/ui/sonner";
import {useRouteTitle} from "@/hooks/use-route-title";
import type React from "react";
import {Outlet} from "react-router-dom";

type AuthLayoutProps = React.ComponentProps<"div">;

const AuthLayout: React.FC<AuthLayoutProps> = (props) => {
  useRouteTitle();

  return (
    <div
      className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10"
      {...props}
    >
      <div className="flex w-full max-w-sm flex-col gap-6">
        <Logo variant="default" />
        <Outlet />
      </div>
      <Toaster />
    </div>
  );
};

export default AuthLayout;
