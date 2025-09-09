import config from "@/lib/config";
import {apiRequest} from "./api.service";
import type {HostResponse, HostsResponse} from "@/types/responses/hosts";
import type {ProgramsResponse} from "@/types/responses/program";

const hostApiURL = `${config.apiUrl}/hosts`;
const hostService = {
  getHosts: () =>
    apiRequest<HostsResponse>({
      protected: true,
      url: hostApiURL,
      method: "GET",
      headers: undefined,
      params: undefined,
    }),
  getHost: (hostId: string) =>
    apiRequest<HostResponse>({
      protected: true,
      url: `${hostApiURL}/${hostId}`,
      method: "GET",
      headers: undefined,
      params: undefined,
    }),
  getHostPrograms: (hostId: string) =>
    apiRequest<ProgramsResponse>({
      protected: true,
      url: `${hostApiURL}/${hostId}/programs`,
      method: "GET",
      headers: undefined,
      params: undefined,
    }),
};

export default hostService;
