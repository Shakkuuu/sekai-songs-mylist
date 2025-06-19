import { useEffect, useState, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { myListClient } from "../lib/grpcClient";
import {
  GetMyListChartByIDRequest,
  AddMyListChartAttachmentRequest,
  DeleteMyListChartAttachmentRequest,
  GetMyListChartAttachmentsByMyListChartIDRequest,
  MyListChart,
  MyListChartAttachment,
} from "../gen/mylist/v1/mylist_pb";
import { AttachmentType, ClearType } from "../gen/enums/mylist_pb";
import { getDifficultyTypeDisplayName } from "../utils/enumDisplay";
import "./MyListChartDetailPage.css";
import { IMAGE_BASE_URL } from "../lib/constants";

const ATTACHMENT_TYPE_LABELS: Record<AttachmentType, string> = {
  [AttachmentType.UNSPECIFIED]: "未指定",
  [AttachmentType.PICTURE]: "画像",
  [AttachmentType.MOVIE]: "動画",
};

const clearTypeOptions = [
  { value: ClearType.UNSPECIFIED, label: "未設定" },
  { value: ClearType.NOT_CLEARED, label: "未クリア" },
  { value: ClearType.CLEARED, label: "クリア" },
  { value: ClearType.FULL_COMBO, label: "フルコンボ" },
  { value: ClearType.ALL_PERFECT, label: "オールパーフェクト" },
];

interface ChartDetailModalProps {
  chart: MyListChart;
  onClose: () => void;
}

const ChartDetailModal = ({ chart, onClose }: ChartDetailModalProps) => {
  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close" onClick={onClose}>
          ×
        </button>
        <h3>{chart.chart?.song?.name}</h3>
        <div className="modal-details">
          <p>
            <strong>難易度:</strong>{" "}
            {getDifficultyTypeDisplayName(chart.chart?.difficultyType ?? 0)}
          </p>
          <p>
            <strong>レベル:</strong> {chart.chart?.level}
          </p>
          <p>
            <strong>譜面ビュー:</strong>{" "}
            <a
              href={chart.chart?.chartViewLink}
              target="_blank"
              rel="noopener noreferrer"
            >
              {chart.chart?.chartViewLink}
            </a>
          </p>
          <p>
            <strong>作詞:</strong> {chart.chart?.song?.lyrics?.name}
          </p>
          <p>
            <strong>作曲:</strong> {chart.chart?.song?.music?.name}
          </p>
          <p>
            <strong>編曲:</strong> {chart.chart?.song?.arrangement?.name}
          </p>
          <p>
            <strong>リリース時間:</strong>{" "}
            {chart.chart?.song?.releaseTime?.toDate().toLocaleString()}
          </p>
          <p>
            <strong>原曲MV:</strong>{" "}
            <a
              href={chart.chart?.song?.originalVideo}
              target="_blank"
              rel="noopener noreferrer"
            >
              {chart.chart?.song?.originalVideo}
            </a>
          </p>
          <p>
            <strong>ボーカルパターン:</strong>
          </p>
          <ul>
            {chart.chart?.song?.vocalPatterns?.map((pattern, index) => (
              <li key={index}>
                {pattern.name} -{" "}
                {pattern.singers?.map((singer) => singer.name).join(", ")}
              </li>
            ))}
          </ul>
          <p>
            <strong>ユニット:</strong>{" "}
            {chart.chart?.song?.units?.map((unit) => unit.name).join(", ")}
          </p>
        </div>
      </div>
    </div>
  );
};

export const MyListChartDetailPage = () => {
  const { myListId, myListChartId } = useParams();
  const [chart, setChart] = useState<MyListChart | null>(null);
  const [clearType, setClearType] = useState<ClearType | undefined>(undefined);
  const [memo, setMemo] = useState("");
  const [attachments, setAttachments] = useState<MyListChartAttachment[]>([]);
  const [attachmentType, setAttachmentType] = useState<AttachmentType>(
    AttachmentType.UNSPECIFIED
  );
  const [fileUrl, setFileUrl] = useState("");
  const [caption, setCaption] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [showChartDetail, setShowChartDetail] = useState(false);
  const navigate = useNavigate();

  const fetchAttachments = useCallback(async () => {
    if (!myListChartId) return;
    const res = await myListClient.getMyListChartAttachmentsByMyListChartID(
      new GetMyListChartAttachmentsByMyListChartIDRequest({
        myListChartId: Number(myListChartId),
      })
    );
    setAttachments(res.myListChartAttachments ?? null);
  }, [myListChartId]);

  useEffect(() => {
    console.log(myListChartId);
    myListClient
      .getMyListChartByID(
        new GetMyListChartByIDRequest({ id: Number(myListChartId) })
      )
      .then((res) => {
        if (res.myListChart) {
          setChart(res.myListChart);
          setClearType(res.myListChart.clearType);
          setMemo(res.myListChart.memo ?? "");
        }
      });
    fetchAttachments();
  }, [myListChartId, fetchAttachments]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selected = e.target.files?.[0];
    setFile(selected ?? null);
    setFileUrl("");
  };

  const handleFileUpload = async () => {
    if (!file) return;
    setUploading(true);
    const formData = new FormData();
    formData.append("file", file);

    try {
      const res = await fetch(IMAGE_BASE_URL + "/upload/attachment", {
        method: "POST",
        body: formData,
      });
      if (!res.ok) throw new Error("アップロードに失敗しました");
      const data = await res.json();
      if (!data.url) throw new Error("サーバーからURLが返されませんでした");
      setFileUrl(data.url);
      alert("ファイルをアップロードしました");
    } catch (err) {
      alert("ファイルのアップロードに失敗しました");
      console.error(err);
    } finally {
      setUploading(false);
    }
  };

  const handleAddAttachment = async () => {
    if (!myListChartId || !fileUrl || !caption) return;
    await myListClient.addMyListChartAttachment(
      new AddMyListChartAttachmentRequest({
        myListChartId: Number(myListChartId),
        attachmentType,
        fileUrl,
        caption,
      })
    );
    await fetchAttachments();
    setFileUrl("");
    setCaption("");
    setAttachmentType(AttachmentType.UNSPECIFIED);
    setFile(null);
  };

  const handleDeleteAttachment = async (id: number) => {
    await myListClient.deleteMyListChartAttachment(
      new DeleteMyListChartAttachmentRequest({ id })
    );
    const att = attachments.find((a) => a.id === id);
    if (att && att.fileUrl) {
      const fileUrl = att.fileUrl.startsWith("http")
        ? new URL(att.fileUrl).pathname
        : att.fileUrl;
      const fileName = fileUrl.split("/").pop();
      if (fileName) {
        fetch(IMAGE_BASE_URL + `/delete/attachment/${fileName}`, {
          method: "DELETE",
        }).catch((err) => {
          console.warn("ファイルサーバーでの削除に失敗:", err);
        });
      }
    }
    await fetchAttachments();
  };

  const renderAttachmentMedia = (att: MyListChartAttachment) => {
    const url = att.fileUrl.startsWith("http")
      ? att.fileUrl
      : IMAGE_BASE_URL + `${att.fileUrl}`;
    if (att.attachmentType === AttachmentType.PICTURE) {
      return (
        <img
          src={url}
          alt={att.caption}
          style={{ maxWidth: 120, maxHeight: 120, display: "block" }}
        />
      );
    }
    if (att.attachmentType === AttachmentType.MOVIE) {
      return (
        <video
          src={url}
          controls
          style={{ maxWidth: 200, maxHeight: 120, display: "block" }}
        >
          お使いのブラウザは動画タグに対応していません
        </video>
      );
    }
    return (
      <a href={url} target="_blank" rel="noopener noreferrer">
        {att.caption}
      </a>
    );
  };

  const handleSave = async () => {
    if (!chart) return;
    await myListClient.changeMyListChartClearType({
      id: chart.id,
      clearType,
    });
    await myListClient.changeMyListChartMemo({
      id: chart.id,
      memo,
    });
    navigate(`/mylist/${myListId}`);
  };

  if (!chart) return null;

  return (
    <div className="mylist-chart-detail-container">
      <div className="mylist-chart-detail-header">
        <h2>譜面詳細</h2>
        <div className="header-actions">
          <button className="save-button" onClick={handleSave}>
            保存
          </button>
          <button
            className="back-button"
            onClick={() => navigate(`/mylist/${myListId}`)}
          >
            戻る
          </button>
        </div>
      </div>
      <div className="chart-detail">
        <div className="chart-thumbnail">
          {chart.chart?.song?.thumbnail ? (
            <img
              src={
                chart.chart.song.thumbnail.startsWith("http")
                  ? chart.chart.song.thumbnail
                  : IMAGE_BASE_URL + `${chart.chart.song.thumbnail}`
              }
              alt={chart.chart.song.name}
            />
          ) : (
            <div className="no-thumbnail">No Image</div>
          )}
        </div>
        <div className="chart-info">
          <div className="chart-header">
            <h3 className="chart-name">{chart.chart?.song?.name}</h3>
            <div className="chart-difficulty">
              {getDifficultyTypeDisplayName(chart.chart?.difficultyType ?? 0)}{" "}
              {chart.chart?.level}
            </div>
          </div>
          <div className="chart-creators">
            <p>作詞: {chart.chart?.song?.lyrics?.name}</p>
            <p>作曲: {chart.chart?.song?.music?.name}</p>
            <p>編曲: {chart.chart?.song?.arrangement?.name}</p>
          </div>
          <button
            className="detail-button"
            onClick={() => setShowChartDetail(true)}
          >
            詳細情報
          </button>
        </div>
      </div>

      <div className="edit-section">
        <h3>編集</h3>
        <div className="edit-fields">
          <div className="status-field">
            <label>クリア状況:</label>
            <select
              value={clearType === undefined ? "" : clearType}
              onChange={(e) => {
                const value = e.target.value;
                setClearType(
                  value === "" ? undefined : (Number(value) as ClearType)
                );
              }}
            >
              {clearTypeOptions.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
          </div>
          <div className="status-field">
            <label>メモ:</label>
            <textarea
              value={memo}
              onChange={(e) => setMemo(e.target.value)}
              placeholder="メモを入力"
            />
          </div>
        </div>
      </div>

      <div className="attachments-section">
        <h3>添付ファイル</h3>
        <ul>
          {attachments.map((att) => (
            <li key={att.id}>
              <span>
                [
                {ATTACHMENT_TYPE_LABELS[att.attachmentType] ??
                  att.attachmentType}
                ]
              </span>{" "}
              {renderAttachmentMedia(att)}
              <span> caption: {att.caption}</span>
              <button
                onClick={() => handleDeleteAttachment(att.id)}
                style={{ marginLeft: 8 }}
              >
                削除
              </button>
            </li>
          ))}
        </ul>
        <div className="attachment-form">
          <label>
            添付種別:{" "}
            <select
              value={attachmentType}
              onChange={(e) => setAttachmentType(Number(e.target.value))}
            >
              {Object.values(AttachmentType)
                .filter((v) => typeof v === "number")
                .map((value) => (
                  <option key={value} value={value}>
                    {ATTACHMENT_TYPE_LABELS[value as AttachmentType]}
                  </option>
                ))}
            </select>
          </label>
          <label style={{ marginLeft: 8 }}>
            ファイル:{" "}
            <input
              type="file"
              accept=".png,.jpeg,.jpg,.mp4,.webm,.mov"
              onChange={handleFileChange}
              style={{ width: 220 }}
            />
          </label>
          <button
            type="button"
            onClick={handleFileUpload}
            disabled={!file || uploading}
            style={{ marginLeft: 8 }}
          >
            {uploading ? "アップロード中..." : "アップロード"}
          </button>
          {fileUrl && (
            <span style={{ marginLeft: 8, color: "green" }}>
              アップロード済み
            </span>
          )}
          <label style={{ marginLeft: 8 }}>
            caption:{" "}
            <input
              value={caption}
              onChange={(e) => setCaption(e.target.value)}
              placeholder="キャプション"
              style={{ width: 120 }}
            />
          </label>
          <button onClick={handleAddAttachment} style={{ marginLeft: 8 }}>
            添付追加
          </button>
        </div>
      </div>

      {showChartDetail && (
        <ChartDetailModal
          chart={chart}
          onClose={() => setShowChartDetail(false)}
        />
      )}
    </div>
  );
};
