import {RouterProvider} from "react-router-dom";
import router from "./router";
import {QueryClientProvider} from "@tanstack/react-query";
import queryClient from "./lib/queryClient";
import useConfigStore from "./stores/config.store";
import {useEffect} from "react";

export default function App() {
  const checkConfiguration = useConfigStore(
    (state) => state.checkConfiguration
  );

  useEffect(() => {
    checkConfiguration();
  }, [checkConfiguration]);
  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  );
}
