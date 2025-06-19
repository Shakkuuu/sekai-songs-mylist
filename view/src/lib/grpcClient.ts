import { createPromiseClient, Interceptor } from "@connectrpc/connect";
import { MasterService } from "../gen/master/master_connect";
import { AuthService } from "../gen/auth/v1/auth_connect";
import { UserService } from "../gen/user/v1/user_connect";
import { MyListService } from "../gen/mylist/v1/mylist_connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { API_BASE_URL } from "./constants";

// エラー型の定義


// 認証ミドルウェア
const authMiddleware: Interceptor = (next) => async (req) => {
  const token: string | null = localStorage.getItem("token");
  if (token) {
    req.header.set("Authorization", `Bearer ${token}`);
  }
  try {
    return await next(req);
  } catch (err: unknown) {
    // 401, 403 など認証エラー時はloginページへ遷移
    if (
      typeof err === "object" &&
      err !== null &&
      "code" in err &&
      (err as { code?: number }).code !== undefined &&
      ((err as { code?: number }).code === 16 ||
        (err as { code?: number }).code === 401 ||
        (err as { code?: number }).code === 403)
    ) {
      window.location.href = "/login";
      throw new Error("Authentication error: redirecting to login");
    }
    throw err;
  }
};

const defaultTransport = createConnectTransport({
  baseUrl: API_BASE_URL,
});

const authTransport = createConnectTransport({
  baseUrl: API_BASE_URL,
  interceptors: [authMiddleware],
});

export const masterClient = createPromiseClient(MasterService, defaultTransport);
export const authClient = createPromiseClient(AuthService, defaultTransport);
export const userClient = createPromiseClient(UserService, authTransport);
export const myListClient = createPromiseClient(MyListService, authTransport);
