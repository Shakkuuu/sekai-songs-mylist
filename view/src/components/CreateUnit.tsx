import { useState } from "react";
import { masterClient } from "../lib/grpcClient";
import { ConnectError } from "@bufbuild/connect";

export const CreateUnit = () => {
  const [name, setName] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await masterClient.createUnit({ name });
      alert("ユニットを作成しました！");
      setName("");
      window.location.reload();
    } catch (error) {
      if (
        error instanceof ConnectError ||
        (error &&
          typeof error === "object" &&
          "name" in error &&
          error.name === "ConnectError")
      ) {
        // gRPC エラーの場合
        const grpcError = error as ConnectError;
        console.error("gRPC Error:", grpcError.code, grpcError.message);
        alert(`Error: ${grpcError.message || String(grpcError)}`);
      } else {
        // その他のエラーの場合
        console.error("Unexpected Error:", error);
        alert("An unexpected error occurred.");
      }
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>ユニット作成</h2>
      <div>
        <label>
          名前:
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </label>
      </div>
      <button type="submit">作成</button>
    </form>
  );
};
