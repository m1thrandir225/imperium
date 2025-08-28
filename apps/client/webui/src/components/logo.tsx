import {Castle} from "lucide-react";
import type React from "react";
import {Link} from "react-router-dom";

type ComponentProps = {
  variant: "default" | "icon";
};

const Logo: React.FC<ComponentProps> = ({variant}) => {
  if (variant === "default") {
    return (
      <Link to="/" className="flex items-center gap-2 self-center font-medium">
        <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
          <Castle className="size-4" />
        </div>
        Imperium
      </Link>
    );
  }

  return (
    <Link to="/" className="flex items-center gap-2 self-center font-medium">
      <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
        <Castle className="size-4" />
      </div>
    </Link>
  );
};

export default Logo;
