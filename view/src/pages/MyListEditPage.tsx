import { useEffect, useState } from "react";
import { masterClient, myListClient } from "../lib/grpcClient";
import { useParams, useNavigate } from "react-router-dom";
import type { Chart } from "../gen/master/chart_pb";
import { AddMyListChartRequest } from "../gen/mylist/v1/mylist_pb";
import { ClearType } from "../gen/enums/mylist_pb";

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
    <div>
      <h2>譜面追加</h2>
      <button onClick={() => navigate(`/mylist/${myListId}`)}>編集終了</button>
      <ul>
        {charts.map((chart) => (
          <li key={chart.id}>
            <div>楽曲名: {chart.song?.name}</div>
            {/* ListChart.tsxと同様の情報を表示 */}
            <button onClick={() => setShowAdd(chart.id)}>追加</button>
            {showAdd === chart.id && (
              <div>
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
                <button onClick={() => handleAdd(chart.id)}>確定</button>
                <button onClick={() => setShowAdd(null)}>キャンセル</button>
              </div>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
};
