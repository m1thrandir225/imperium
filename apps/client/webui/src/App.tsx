import {QueryClientProvider} from "@tanstack/react-query";
import {useEffect} from "react";
import {RouterProvider} from "react-router-dom";
import queryClient from "./lib/queryClient";
import router from "./router";
import useConfigStore from "./stores/config.store";

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
