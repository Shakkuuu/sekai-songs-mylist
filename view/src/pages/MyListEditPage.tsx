import { useEffect, useState } from "react";
import { masterClient, myListClient } from "../lib/grpcClient";
import { useParams, useNavigate } from "react-router-dom";
import type { Chart } from "../gen/master/chart_pb";
import { AddMyListChartRequest } from "../gen/mylist/v1/mylist_pb";
import { ClearType } from "../gen/enums/mylist_pb";
import { getDifficultyTypeDisplayName } from "../utils/enumDisplay";
import "./MyListEditPage.css";

const clearTypeOptions = [
  { value: ClearType.UNSPECIFIED, label: "未設定" },
  { value: ClearType.NOT_CLEARED, label: "未クリア" },
  { value: ClearType.CLEARED, label: "クリア" },
  { value: ClearType.FULL_COMBO, label: "フルコンボ" },
  { value: ClearType.ALL_PERFECT, label: "オールパーフェクト" },
];

export const MyListEditPage = () => {
  const { myListId } = useParams();
  const [charts, setCharts] = useState<Chart[]>([]);
  const [showAdd, setShowAdd] = useState<number | null>(null);
  const [clearType, setClearType] = useState<ClearType | undefined>(undefined);
  const [memo, setMemo] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    masterClient.getCharts({}).then((res) => setCharts(res.charts));
  }, []);

  const handleAdd = async (chartId: number) => {
    await myListClient.addMyListChart(
      new AddMyListChartRequest({
        myListId: Number(myListId),
        chartId,
        clearType,
        memo,
      })
    );
    setShowAdd(null);
    setClearType(undefined);
    setMemo("");
    // 追加後は編集画面に戻る
    navigate(`/mylist/${myListId}`);
  };

  return (
    <div className="mylist-edit-container">
      <div className="mylist-edit-header">
        <h2>譜面追加</h2>
        <button
          className="back-button"
          onClick={() => navigate(`/mylist/${myListId}`)}
        >
          編集終了
        </button>
      </div>
      <div className="chart-list">
        {charts.map((chart) => (
          <div key={chart.id} className="chart-item">
            <div className="chart-thumbnail">
              {chart.song?.thumbnail &&
              chart.song.thumbnail.startsWith("http") ? (
                <img src={chart.song.thumbnail} alt={chart.song.name} />
              ) : (
                <div className="no-thumbnail">No Image</div>
              )}
            </div>
            <div className="chart-info">
              <div className="chart-header">
                <h3 className="chart-name">{chart.song?.name}</h3>
                <div className="chart-difficulty">
                  {getDifficultyTypeDisplayName(chart.difficultyType)}{" "}
                  {chart.level}
                </div>
              </div>
              <p className="chart-creators">
                {chart.song?.lyrics?.name} / {chart.song?.music?.name} /{" "}
                {chart.song?.arrangement?.name}
              </p>
            </div>
            <div className="chart-actions">
              <button
                className="add-button"
                onClick={() => setShowAdd(chart.id)}
              >
                追加
              </button>
              {showAdd === chart.id && (
                <div className="add-form">
                  <select
                    value={clearType === undefined ? "" : clearType}
                    onChange={(e) => {
                      const value = e.target.value;
                      setClearType(
                        value === "" ? undefined : (Number(value) as ClearType)
                      );
                    }}
                  >
                    <option value="">クリア状況を選択</option>
                    {clearTypeOptions.map((opt) => (
                      <option key={opt.value} value={opt.value}>
                        {opt.label}
                      </option>
                    ))}
                  </select>
                  <input
                    value={memo}
                    onChange={(e) => setMemo(e.target.value)}
                    placeholder="メモ"
                  />
                  <div className="add-form-buttons">
                    <button onClick={() => handleAdd(chart.id)}>確定</button>
                    <button onClick={() => setShowAdd(null)}>キャンセル</button>
                  </div>
                </div>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
