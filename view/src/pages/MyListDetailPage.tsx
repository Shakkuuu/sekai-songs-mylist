import { useEffect, useState } from "react";
import { myListClient } from "../lib/grpcClient";
import { useParams, useNavigate } from "react-router-dom";
import {
  MyListChart,
  GetMyListChartsByMyListIDRequest,
  DeleteMyListChartRequest,
} from "../gen/mylist/v1/mylist_pb";
import {
  getDifficultyTypeDisplayName,
  getClearTypeDisplayName,
} from "../utils/enumDisplay";

export const MyListDetailPage = () => {
  const { myListId } = useParams();
  const [charts, setCharts] = useState<MyListChart[]>([]);
  const navigate = useNavigate();

  useEffect(() => {
    myListClient
      .getMyListChartsByMyListID(
        new GetMyListChartsByMyListIDRequest({ myListId: Number(myListId) })
      )
      .then((res) => setCharts(res.myListCharts));
  }, [myListId]);

  return (
    <div>
      <h2>マイリスト詳細</h2>
      <button onClick={() => navigate(`/mylist/${myListId}/edit`)}>編集</button>
      <ul>
        {charts
          .slice()
          .sort((a, b) => (a.chart?.id ?? 0) - (b.chart?.id ?? 0))
          .map((mlc) => (
            <li key={mlc.id}>
              <div>楽曲名: {mlc.chart?.song?.name}</div>
              <div>作詞: {mlc.chart?.song?.lyrics?.name}</div>
              <div>作曲: {mlc.chart?.song?.music?.name}</div>
              <div>編曲: {mlc.chart?.song?.arrangement?.name}</div>
              <div>
                歌唱:{" "}
                {mlc.chart?.song?.vocalPatterns
                  ?.map((vp) => vp.singers?.map((s) => s.name).join(", "))
                  .join(" / ")}
              </div>
              <div>
                譜面:{" "}
                {mlc.chart?.difficultyType !== undefined
                  ? getDifficultyTypeDisplayName(mlc.chart.difficultyType)
                  : ""}{" "}
                {mlc.chart?.level}
              </div>
              <div>クリア状況: {getClearTypeDisplayName(mlc.clearType)}</div>
              <div>メモ: {mlc.memo}</div>
              <button
                onClick={() => navigate(`/mylist/${myListId}/chart/${mlc.id}`)}
              >
                編集
              </button>
              <button
                onClick={async () => {
                  await myListClient.deleteMyListChart(
                    new DeleteMyListChartRequest({ id: mlc.id })
                  );
                  setCharts(charts.filter((c) => c.id !== mlc.id));
                }}
              >
                削除
              </button>
            </li>
          ))}
      </ul>
    </div>
  );
};
