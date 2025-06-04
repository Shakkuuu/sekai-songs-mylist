import { useState } from "react";
import { authClient } from "../lib/grpcClient";
import { useNavigate } from "react-router-dom";

export const SignupPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [checkPassword, setCheckPassword] = useState("");
  const [message, setMessage] = useState("");
  const navigate = useNavigate();

  const handleSignup = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await authClient.signup({ email, password, checkPassword });
      setMessage("サインアップ成功！ログイン画面に移行します。");
      setTimeout(() => navigate("/login"), 1000);
    } catch (error) {
      setMessage(
        "サインアップ失敗: " +
          (error && typeof error === "object" && "message" in error
            ? (error as { message: string }).message
            : String(error))
      );
    }
  };

  return (
    <div>
      <h2>Signup</h2>
      <form onSubmit={handleSignup}>
        <div>
          <label>
            Email:{" "}
            <input
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </label>
        </div>
        <div>
          <label>
            Password:{" "}
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </label>
        </div>
        <div>
          <label>
            パスワードの確認:{" "}
            <input
              type="password"
              value={checkPassword}
              onChange={(e) => setCheckPassword(e.target.value)}
              required
            />
          </label>
        </div>
        <button type="submit">Signup</button>
      </form>
      {message && <div>{message}</div>}
    </div>
  );
};
