import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { myListClient } from "../lib/grpcClient";
import {
  GetMyListChartByIDRequest,
  ChangeMyListChartClearTypeRequest,
  ChangeMyListChartMemoRequest,
  AddMyListChartAttachmentRequest,
  DeleteMyListChartAttachmentRequest,
  GetMyListChartAttachmentsByMyListChartIDRequest,
  MyListChart,
  MyListChartAttachment,
} from "../gen/mylist/v1/mylist_pb";
import { AttachmentType, ClearType } from "../gen/enums/mylist_pb";
import { getClearTypeDisplayName } from "../utils/enumDisplay";

const ATTACHMENT_TYPE_LABELS: Record<AttachmentType, string> = {
  [AttachmentType.UNSPECIFIED]: "未指定",
  [AttachmentType.PICTURE]: "画像",
  [AttachmentType.MOVIE]: "動画",
};

const CLEAR_TYPE_OPTIONS = [
  { value: ClearType.UNSPECIFIED, label: "未設定" },
  { value: ClearType.NOT_CLEARED, label: "未クリア" },
  { value: ClearType.CLEARED, label: "クリア" },
  { value: ClearType.FULL_COMBO, label: "フルコンボ" },
  { value: ClearType.ALL_PERFECT, label: "オールパーフェクト" },
];

export const MyListChartDetailPage = () => {
  const { myListId, myListChartId } = useParams();
  const navigate = useNavigate();
  const [myListChart, setMyListChart] = useState<MyListChart | null>(null);
  const [clearType, setClearType] = useState<ClearType>(ClearType.UNSPECIFIED);
  const [memo, setMemo] = useState("");
  const [attachments, setAttachments] = useState<MyListChartAttachment[]>([]);
  const [attachmentType, setAttachmentType] = useState<AttachmentType>(
    AttachmentType.UNSPECIFIED
  );
  const [fileUrl, setFileUrl] = useState("");
  const [caption, setCaption] = useState("");

  // MyListChart取得
  useEffect(() => {
    if (!myListChartId) return;
    myListClient
      .getMyListChartByID(
        new GetMyListChartByIDRequest({ id: Number(myListChartId) })
      )
      .then((res) => {
        setMyListChart(res.myListChart ?? null);
        setClearType(res.myListChart?.clearType ?? ClearType.UNSPECIFIED);
        setMemo(res.myListChart?.memo ?? "");
      });
    fetchAttachments();
    // eslint-disable-next-line
  }, [myListChartId]);

  // 添付ファイル一覧取得
  const fetchAttachments = async () => {
    if (!myListChartId) return;
    const res = await myListClient.getMyListChartAttachmentsByMyListChartID(
      new GetMyListChartAttachmentsByMyListChartIDRequest({
        myListChartId: Number(myListChartId),
      })
    );
    setAttachments(res.myListChartAttachments ?? null);
  };

  // クリア状況更新
  const handleClearTypeUpdate = async () => {
    if (!myListChartId) return;
    await myListClient.changeMyListChartClearType(
      new ChangeMyListChartClearTypeRequest({
        id: Number(myListChartId),
        clearType,
      })
    );
    const res = await myListClient.getMyListChartByID(
      new GetMyListChartByIDRequest({ id: Number(myListChartId) })
    );
    setMyListChart(res.myListChart ?? null);
    setClearType(res.myListChart?.clearType ?? ClearType.UNSPECIFIED);
  };

  // メモ更新
  const handleMemoUpdate = async () => {
    if (!myListChartId) return;
    await myListClient.changeMyListChartMemo(
      new ChangeMyListChartMemoRequest({
        id: Number(myListChartId),
        memo,
      })
    );
    const res = await myListClient.getMyListChartByID(
      new GetMyListChartByIDRequest({ id: Number(myListChartId) })
    );
    setMyListChart(res.myListChart ?? null);
    setMemo(res.myListChart ? res.myListChart.memo ?? "" : "");
  };

  // 添付追加
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
  };

  // 添付削除
  const handleDeleteAttachment = async (id: number) => {
    await myListClient.deleteMyListChartAttachment(
      new DeleteMyListChartAttachmentRequest({ id })
    );
    await fetchAttachments();
  };

  return (
    <div>
      <h2>MyListChart詳細・編集</h2>
      <button onClick={() => navigate(`/mylist/${myListId}`)}>戻る</button>
      {myListChart && (
        <div style={{ margin: "1em 0", background: "#f5f5f5", padding: "1em" }}>
          <div>
            <strong>楽曲名:</strong> {myListChart.chart?.song?.name}
          </div>
          <div>
            <strong>難易度:</strong> {myListChart.chart?.difficultyType}
          </div>
          <div>
            <strong>レベル:</strong> {myListChart.chart?.level}
          </div>
        </div>
      )}
      <div>
        <label>
          クリア状況:{" "}
          <select
            value={clearType}
            onChange={(e) => setClearType(Number(e.target.value))}
          >
            {CLEAR_TYPE_OPTIONS.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </label>
        <button onClick={handleClearTypeUpdate} style={{ marginLeft: 8 }}>
          クリア状況を保存
        </button>
        <span style={{ marginLeft: 8 }}>
          {getClearTypeDisplayName(clearType)}
        </span>
      </div>
      <div style={{ marginTop: 12 }}>
        <label>
          メモ:{" "}
          <input
            value={memo}
            onChange={(e) => setMemo(e.target.value)}
            style={{ width: 200 }}
          />
        </label>
        <button onClick={handleMemoUpdate} style={{ marginLeft: 8 }}>
          メモを保存
        </button>
      </div>
      <h3 style={{ marginTop: 24 }}>添付ファイル</h3>
      <ul>
        {attachments.map((att) => (
          <li key={att.id}>
            <span>
              [
              {ATTACHMENT_TYPE_LABELS[att.attachmentType] ?? att.attachmentType}
              ]
            </span>{" "}
            <a href={att.fileUrl} target="_blank" rel="noopener noreferrer">
              {att.caption}
            </a>{" "}
            caption: {att.caption}
            <button
              onClick={() => handleDeleteAttachment(att.id)}
              style={{ marginLeft: 8 }}
            >
              削除
            </button>
          </li>
        ))}
      </ul>
      <div style={{ marginTop: 12 }}>
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
          file_url:{" "}
          <input
            value={fileUrl}
            onChange={(e) => setFileUrl(e.target.value)}
            placeholder="https://example.com/file.png"
            style={{ width: 220 }}
          />
        </label>
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
  );
};
