import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { userClient } from "../lib/grpcClient";

export function useAdminGuard() {
  const navigate = useNavigate();
  const [checked, setChecked] = useState(false);

  useEffect(() => {
    async function check() {
      try {
        // IsAdminリクエストは空でOK
        const res = await userClient.isAdmin({});
        if (!res.isAdmin) {
          navigate("/login", { replace: true });
        } else {
          setChecked(true);
        }
      } catch {
        // 認証エラー時も/loginへ
        navigate("/login", { replace: true });
      }
    }
    check();
  }, [navigate]);

  return checked;
}
