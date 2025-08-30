import config from "@/lib/config";
import {apiRequest} from "./api.service";
import type {
  ConfigStatus,
  SetupConfigRequest,
  SetupConfigResponse,
} from "@/types/responses/config";

const configURL = `${config.apiUrl}/config`;

const configService = {
  getConfigStatus: () =>
    apiRequest<ConfigStatus>({
      url: `${configURL}/status`,
      method: "GET",
      params: undefined,
      protected: false,
      headers: undefined,
    }),

  setupConfig: (input: SetupConfigRequest) =>
    apiRequest<SetupConfigResponse>({
      url: `${configURL}/setup`,
      method: "POST",
      data: input,
      params: undefined,
      protected: false,
      headers: undefined,
    }),
};

export default configService;
