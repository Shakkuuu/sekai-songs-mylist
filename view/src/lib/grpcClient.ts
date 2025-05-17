import { createPromiseClient } from "@connectrpc/connect";
import { MasterService } from "../gen/master/master_connect";
import { createConnectTransport } from "@connectrpc/connect-web";

const transport = createConnectTransport({
  baseUrl: "http://localhost:8080", // gRPC Gateway の URL に合わせて
});

export const masterClient = createPromiseClient(MasterService, transport);
