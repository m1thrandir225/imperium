import {useEffect} from "react";
import {useMatches} from "react-router-dom";

type RouteHandle = {
  title?: () => string;
};

export const useRouteTitle = () => {
  const matches = useMatches();

  useEffect(() => {
    const match = [...matches]
      .reverse()
      .find((match) => (match.handle as RouteHandle)?.title);

    if (match?.handle && (match.handle as RouteHandle).title) {
      document.title = (match.handle as RouteHandle).title!();
    } else {
      document.title = "Imperium";
    }
  }, [matches]);
};
