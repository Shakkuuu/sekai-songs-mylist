import { useState } from "react";
import { masterClient } from "../lib/grpcClient";
import { ConnectError } from "@bufbuild/connect";
import { useSongs, useSingers } from "../hooks/useMasterLists";
import Select from "react-select";

export const CreateVocalPattern = () => {
  const songs = useSongs();
  const singers = useSingers();

  const [songId, setSongId] = useState<number>(0);
  const [name, setName] = useState<string>("");
  const [selectedSingerIds, setSelectedSingerIds] = useState<number[]>([]);
  const [singerPositions, setSingerPositions] = useState<string>(""); // comma-separated positions

  const singerOptions = singers.map((s) => ({ value: s.id, label: s.name }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const singerPositionsArray = singerPositions
      .split(",")
      .map((pos) => parseInt(pos.trim(), 10));

    try {
      await masterClient.createVocalPattern({
        songId,
        name,
        singerIds: selectedSingerIds,
        singerPositions: singerPositionsArray,
      });
      alert("VocalPattern created successfully!");
      setSongId(0);
      setName("");
      setName("");
      setSelectedSingerIds([]);
      setSingerPositions("");
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
        if ("code" in error && "message" in error) {
          const grpcError = error as ConnectError;
          console.error("gRPC Error:", grpcError.code, grpcError.message);
          alert(`Error: ${grpcError.message || String(error)}`);
        } else {
          console.error("gRPC Error:", error);
          alert(`Error: ${String(error)}`);
        }
      } else {
        // その他のエラーの場合
        console.error("Unexpected Error:", error);
        alert("An unexpected error occurred.");
      }
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>Create VocalPattern</h2>
      <div>
        <label>
          Song:
          <select
            value={songId}
            onChange={(e) => setSongId(Number(e.target.value))}
            required
          >
            <option value={0}>選択してください</option>
            {songs.map((song) => (
              <option key={song.id} value={song.id}>
                {song.name}
              </option>
            ))}
          </select>
        </label>
      </div>
      <div>
        <label>
          Name:
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </label>
      </div>
      <div>
        <label>
          Singers:
          <Select
            isMulti
            options={singerOptions}
            value={singerOptions.filter((opt) =>
              selectedSingerIds.includes(opt.value)
            )}
            onChange={(opts) =>
              setSelectedSingerIds(opts.map((opt) => opt.value))
            }
          />
        </label>
      </div>
      <div>
        <label>
          Singer Positions (comma-separated):
          <input
            type="text"
            value={singerPositions}
            onChange={(e) => setSingerPositions(e.target.value)}
          />
        </label>
      </div>
      <button type="submit">Create</button>
    </form>
  );
};
