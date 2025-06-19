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
import "./MyListDetailPage.css";
import { IMAGE_BASE_URL } from "../lib/constants";

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
    <div className="mylist-detail-container">
      <div className="mylist-detail-header">
        <h2>マイリスト詳細</h2>
        <div className="header-actions">
          <button
            className="edit-button"
            onClick={() => navigate(`/mylist/${myListId}/edit`)}
          >
            編集
          </button>
          <button className="back-button" onClick={() => navigate("/mylist")}>
            戻る
          </button>
        </div>
      </div>
      <div className="chart-list">
        {charts
          .slice()
          .sort((a, b) => (a.chart?.id ?? 0) - (b.chart?.id ?? 0))
          .map((mlc) => (
            <div key={mlc.id} className="chart-item">
              <div className="chart-thumbnail">
                {mlc.chart?.song?.thumbnail ? (
                  <img
                    src={
                      mlc.chart.song.thumbnail.startsWith("http")
                        ? mlc.chart.song.thumbnail
                        : `${IMAGE_BASE_URL}${mlc.chart.song.thumbnail}`
                    }
                    alt={mlc.chart.song.name}
                  />
                ) : (
                  <div className="no-thumbnail">No Image</div>
                )}
              </div>
              <div className="chart-info">
                <div className="chart-header">
                  <h3 className="chart-name">{mlc.chart?.song?.name}</h3>
                  <div className="chart-difficulty">
                    {getDifficultyTypeDisplayName(
                      mlc.chart?.difficultyType ?? 0
                    )}{" "}
                    {mlc.chart?.level}
                  </div>
                </div>
                <p className="chart-creators">
                  {mlc.chart?.song?.lyrics?.name} /{" "}
                  {mlc.chart?.song?.music?.name} /{" "}
                  {mlc.chart?.song?.arrangement?.name}
                </p>
                <div className="chart-status">
                  <span className="clear-status">
                    クリア状況: {getClearTypeDisplayName(mlc.clearType)}
                  </span>
                  {mlc.memo && <span className="memo">メモ: {mlc.memo}</span>}
                </div>
              </div>
              <div className="chart-actions">
                <button
                  className="edit-button"
                  onClick={() =>
                    navigate(`/mylist/${myListId}/chart/${mlc.id}`)
                  }
                >
                  詳細
                </button>
                <button
                  className="delete-button"
                  onClick={async () => {
                    await myListClient.deleteMyListChart(
                      new DeleteMyListChartRequest({ id: mlc.id })
                    );
                    setCharts(charts.filter((c) => c.id !== mlc.id));
                  }}
                >
                  削除
                </button>
              </div>
            </div>
          ))}
      </div>
    </div>
  );
};
