import { useState } from "react";
import { masterClient } from "../lib/grpcClient";
import { Timestamp } from "@bufbuild/protobuf"; // Timestamp をインポート
import { ConnectError } from "@bufbuild/connect"; // ConnectError をインポート
import { MusicVideoType } from "../gen/enums/master_pb";
import { useArtists, useUnits } from "../hooks/useMasterLists";
import Select from "react-select";

export const CreateSong = () => {
  const artists = useArtists();
  const units = useUnits();

  const [name, setName] = useState<string>("");
  const [kana, setKana] = useState<string>("");
  const [lyricsId, setLyricsId] = useState<number>(0);
  const [musicId, setMusicId] = useState<number>(0);
  const [arrangementId, setArrangementId] = useState<number>(0);
  const [thumbnail, setThumbnail] = useState<string>("");
  const [thumbnailFile, setThumbnailFile] = useState<File | null>(null); // ファイルを保持
  const [originalVideo, setOriginalVideo] = useState<string>("");
  const [releaseTime, setReleaseTime] = useState<string>(""); // ISO string
  const [deleted, setDeleted] = useState<boolean>(false);
  const [selectedUnitIds, setSelectedUnitIds] = useState<number[]>([]);
  const [musicVideoTypes, setMusicVideoTypes] = useState<MusicVideoType[]>([]);
  const [uploading, setUploading] = useState<boolean>(false);

  const unitOptions = units.map((u) => ({ value: u.id, label: u.name }));

  const handleMusicVideoTypeChange = (type: MusicVideoType) => {
    setMusicVideoTypes((prev) =>
      prev.includes(type) ? prev.filter((t) => t !== type) : [...prev, type]
    );
  };

  // サムネイルファイル選択時の処理
  const handleThumbnailChange = async (
    e: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = e.target.files?.[0];
    if (!file) return;
    // 拡張子チェック
    if (!["image/png", "image/jpeg", "image/jpg"].includes(file.type)) {
      alert("png, jpeg, jpg のみアップロードできます");
      return;
    }
    setThumbnailFile(file);
    setThumbnail(""); // URLは一旦クリア
  };

  // サムネイルアップロード処理
  const handleThumbnailUpload = async () => {
    if (!thumbnailFile) return;
    setUploading(true);
    const formData = new FormData();
    formData.append("file", thumbnailFile);

    try {
      const res = await fetch("http://localhost:8888/upload/thumbnail", {
        method: "POST",
        body: formData,
      });
      if (!res.ok) {
        throw new Error("アップロードに失敗しました");
      }
      const data = await res.json();
      if (!data.url) {
        throw new Error("サーバーからURLが返されませんでした");
      }
      setThumbnail(data.url);
      alert("サムネイルをアップロードしました");
    } catch (err) {
      alert("サムネイルのアップロードに失敗しました");
      console.error(err);
    } finally {
      setUploading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // サムネイルが未アップロードの場合は警告
    if (!thumbnail) {
      alert("サムネイル画像をアップロードしてください");
      return;
    }

    console.log(Timestamp.fromDate(new Date(releaseTime)));
    try {
      await masterClient.createSong({
        name,
        kana,
        lyricsId,
        musicId,
        arrangementId,
        thumbnail,
        originalVideo,
        releaseTime: Timestamp.fromDate(new Date(releaseTime)),
        deleted,
        unitIds: selectedUnitIds,
        musicVideoTypes,
      });

      alert("Song created successfully!");
      setName("");
      setKana("");
      setLyricsId(0);
      setMusicId(0);
      setArrangementId(0);
      setThumbnail("");
      setThumbnailFile(null);
      setOriginalVideo("");
      setReleaseTime("");
      setDeleted(false);
      setSelectedUnitIds([]);
      setMusicVideoTypes([]);
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
        if ("code" in error) {
          if (
            typeof error === "object" &&
            error !== null &&
            "code" in error &&
            "message" in error
          ) {
            console.error(
              "gRPC Error:",
              (error as { code: unknown; message: unknown }).code,
              (error as { code: unknown; message: unknown }).message
            );
          } else {
            console.error("gRPC Error:", error);
          }
        } else {
          console.error("gRPC Error:", (error as ConnectError).message);
        }
        alert(
          `Error: ${
            typeof error === "object" && error !== null && "message" in error
              ? (error as { message: string }).message
              : String(error)
          }`
        );
      } else {
        // その他のエラーの場合
        console.error("Unexpected Error:", error);
        alert("An unexpected error occurred.");
      }
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>Create Song</h2>
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
          Kana:
          <input
            type="text"
            value={kana}
            onChange={(e) => setKana(e.target.value)}
            required
          />
        </label>
      </div>
      <div>
        <label>
          Lyrics:
          <select
            value={lyricsId}
            onChange={(e) => setLyricsId(Number(e.target.value))}
            required
          >
            <option value={0}>選択してください</option>
            {artists.map((artist) => (
              <option key={artist.id} value={artist.id}>
                {artist.name}
              </option>
            ))}
          </select>
        </label>
      </div>
      <div>
        <label>
          Music:
          <select
            value={musicId}
            onChange={(e) => setMusicId(Number(e.target.value))}
            required
          >
            <option value={0}>選択してください</option>
            {artists.map((artist) => (
              <option key={artist.id} value={artist.id}>
                {artist.name}
              </option>
            ))}
          </select>
        </label>
      </div>
      <div>
        <label>
          Arrangement:
          <select
            value={arrangementId}
            onChange={(e) => setArrangementId(Number(e.target.value))}
            required
          >
            <option value={0}>選択してください</option>
            {artists.map((artist) => (
              <option key={artist.id} value={artist.id}>
                {artist.name}
              </option>
            ))}
          </select>
        </label>
      </div>
      <div>
        <label>
          Thumbnail:
          <input
            type="file"
            accept=".png,.jpeg,.jpg"
            onChange={handleThumbnailChange}
          />
        </label>
        <button
          type="button"
          onClick={handleThumbnailUpload}
          disabled={!thumbnailFile || uploading}
          style={{ marginLeft: 8 }}
        >
          {uploading ? "アップロード中..." : "アップロード"}
        </button>
        {thumbnail && (
          <span style={{ marginLeft: 8 }}>
            <a href={thumbnail} target="_blank" rel="noopener noreferrer">
              サムネイル画像を確認
            </a>
          </span>
        )}
      </div>
      <div>
        <label>
          Original Video:
          <input
            type="text"
            value={originalVideo}
            onChange={(e) => setOriginalVideo(e.target.value)}
          />
        </label>
      </div>
      <div>
        <label>
          Release Time (ISO format):
          <input
            type="datetime-local"
            value={releaseTime}
            onChange={(e) => setReleaseTime(e.target.value)}
          />
        </label>
      </div>
      <div>
        <label>
          Deleted:
          <input
            type="checkbox"
            checked={deleted}
            onChange={(e) => setDeleted(e.target.checked)}
          />
        </label>
      </div>
      <div>
        <label>
          Units:
          <Select
            isMulti
            options={unitOptions}
            value={unitOptions.filter((opt) =>
              selectedUnitIds.includes(opt.value)
            )}
            onChange={(opts) =>
              setSelectedUnitIds(opts.map((opt) => opt.value))
            }
          />
        </label>
      </div>
      <div>
        <label>
          MusicVideoTypes:
          <div>
            <label>
              <input
                type="checkbox"
                checked={musicVideoTypes.includes(
                  MusicVideoType.MUSIC_VIDEO_TYPE_3D
                )}
                onChange={() =>
                  handleMusicVideoTypeChange(MusicVideoType.MUSIC_VIDEO_TYPE_3D)
                }
              />
              3DMV
            </label>
            <label>
              <input
                type="checkbox"
                checked={musicVideoTypes.includes(
                  MusicVideoType.MUSIC_VIDEO_TYPE_2D
                )}
                onChange={() =>
                  handleMusicVideoTypeChange(MusicVideoType.MUSIC_VIDEO_TYPE_2D)
                }
              />
              2DMV
            </label>
            <label>
              <input
                type="checkbox"
                checked={musicVideoTypes.includes(
                  MusicVideoType.MUSIC_VIDEO_TYPE_ORIGINAL
                )}
                onChange={() =>
                  handleMusicVideoTypeChange(
                    MusicVideoType.MUSIC_VIDEO_TYPE_ORIGINAL
                  )
                }
              />
              原曲MV
            </label>
          </div>
        </label>
      </div>
      <button type="submit">Create</button>
    </form>
  );
};
