import {AppContent} from "@/components/app-content";
import AppHeader from "@/components/app-header";
import type React from "react";

type DefaultLayoutProps = React.ComponentProps<"div">;

const DefaultLayout: React.FC<DefaultLayoutProps> = ({children, ...props}) => {
  return (
    <div className="flex flex-col h-screen gap-4">
      <AppHeader />
      <AppContent {...props}>{children}</AppContent>
    </div>
  );
};

export default DefaultLayout;
