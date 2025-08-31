import {AppContent} from "@/components/app-content";
import AppHeader from "@/components/app-header";
import {Toaster} from "@/components/ui/sonner";
import {useRouteTitle} from "@/hooks/use-route-title";
import type React from "react";
import {Outlet} from "react-router-dom";

type DefaultLayoutProps = React.ComponentProps<"div">;

const DefaultLayout: React.FC<DefaultLayoutProps> = (props) => {
  useRouteTitle();

  return (
    <div className="flex flex-col h-screen gap-4" {...props}>
      <AppHeader />
      <AppContent>
        <Outlet />
      </AppContent>
      <Toaster />
    </div>
  );
};

export default DefaultLayout;
