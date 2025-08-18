import {Link} from "react-router-dom";
import {NavUser} from "./nav-user";

type AppHeaderProps = React.ComponentProps<"div">;

const AppHeader: React.FC<AppHeaderProps> = () => {
  return (
    <>
      <div className="border-sidebar-border/80 border-b">
        <div className="mx-auto flex h-16 items-center px-4 md:max-w-7xl">
          <Link to={"/"} className="flex items-center space-x-2">
            Imperium
          </Link>

          <div className="ml-auto flex items-center space-x-2">
            <NavUser
              user={{
                name: "John Doe",
                email: "john.doe@example.com",
                avatar: "https://github.com/shadcn.png",
              }}
            />
          </div>
        </div>
      </div>
    </>
  );
};

export default AppHeader;
