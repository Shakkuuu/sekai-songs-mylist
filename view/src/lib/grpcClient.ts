import { createPromiseClient } from "@connectrpc/connect";
import { MasterService } from "../gen/master/master_connect";
import { AuthService } from "../gen/auth/v1/auth_connect";
import { UserService } from "../gen/user/v1/user_connect";
import { createConnectTransport } from "@connectrpc/connect-web";

const authMiddleware: PromiseClientMiddleware = (next) => async (req) => {
  const token = localStorage.getItem("token");
  if (token) {
    req.header.set("Authorization", `Bearer ${token}`);
  }
  return next(req);
};

const defaultTransport = createConnectTransport({
  baseUrl: "http://localhost:8080", // gRPC Gateway の URL に合わせて
});

const authTransport = createConnectTransport({
  baseUrl: "http://localhost:8080", // gRPC Gateway の URL に合わせて
  interceptors: [authMiddleware],
});

export const masterClient = createPromiseClient(MasterService, defaultTransport);
export const authClient = createPromiseClient(AuthService, defaultTransport);
export const userClient = createPromiseClient(UserService, authTransport);
