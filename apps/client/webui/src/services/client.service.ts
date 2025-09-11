import config from "@/lib/config";
import {apiRequest} from "./api.service";
import type {Client} from "@/types/models/client";

const clientsAPIUrl = `${config.apiUrl}/clients`;

const clientService = {
  getClientInfo: () =>
    apiRequest<Client>({
      method: "GET",
      url: clientsAPIUrl,
      data: undefined,
      protected: true,
      headers: undefined,
      params: undefined,
    }),
};

export default clientService;
