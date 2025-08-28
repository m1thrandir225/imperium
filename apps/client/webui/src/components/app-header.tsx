import Logo from "./logo";
import {NavUser} from "./nav-user";

type AppHeaderProps = React.ComponentProps<"div">;

const AppHeader: React.FC<AppHeaderProps> = () => {
  return (
    <>
      <div className="border-sidebar-border/80 border-b">
        <div className="mx-auto flex h-16 items-center px-4 md:max-w-7xl">
          <Logo variant="default" />

          <div className="ml-auto flex items-center space-x-2">
            <NavUser />
          </div>
        </div>
      </div>
    </>
  );
};

export default AppHeader;
