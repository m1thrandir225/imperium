import {Castle} from "lucide-react";
import type React from "react";
import {Link} from "react-router-dom";

type AuthLayoutProps = React.ComponentProps<"div">;

const AuthLayout: React.FC<AuthLayoutProps> = ({children, ...props}) => {
  return (
    <div
      {...props}
      className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10"
    >
      <div className="flex w-full max-w-sm flex-col gap-6">
        <Link
          to="/"
          className="flex items-center gap-2 self-center font-medium"
        >
          <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
            <Castle className="size-4" />
          </div>
          Imperium
        </Link>
        {children}
      </div>
    </div>
  );
};

export default AuthLayout;
