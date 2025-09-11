import clientService from "@/services/client.service";
import useAuthStore from "@/stores/auth.store";
import {useQuery} from "@tanstack/react-query";

export function useClientInfo() {
  const isAuthenticated = useAuthStore(
    (state) => !!state.user && !!state.accessToken
  );
  const setClient = useAuthStore((state) => state.setClient);

  const {
    data: client,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["client-info"],
    queryFn: async () => {
      const clientData = await clientService.getClientInfo();
      setClient(clientData);
      return clientData;
    },
    enabled: isAuthenticated,
    staleTime: 1000 * 60 * 5, // 5minutes
    retry: 1,
    refetchOnWindowFocus: false,
  });

  return {
    client,
    isLoading,
    error,
    clientId: client?.id,
    clientName: client?.client_name,
    clientIp: client?.ip_address,
  };
}
