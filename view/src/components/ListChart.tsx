import { useEffect, useState } from "react";
import { masterClient } from "../lib/grpcClient";
import { Chart } from "../gen/master/chart_pb";
import { getDifficultyTypeDisplayName } from "../utils/enumDisplay";
import "./ListChart.css";

interface ChartDetailModalProps {
  chart: Chart;
  onClose: () => void;
}

const ChartDetailModal = ({ chart, onClose }: ChartDetailModalProps) => {
  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close" onClick={onClose}>×</button>
        <h3>{chart.song?.name}</h3>
        <div className="modal-details">
          <p><strong>難易度:</strong> {getDifficultyTypeDisplayName(chart.difficultyType)}</p>
          <p><strong>レベル:</strong> {chart.level}</p>
          <p><strong>譜面ビュー:</strong> <a href={chart.chartViewLink} target="_blank" rel="noopener noreferrer">{chart.chartViewLink}</a></p>
          <p><strong>作詞:</strong> {chart.song?.lyrics?.name}</p>
          <p><strong>作曲:</strong> {chart.song?.music?.name}</p>
          <p><strong>編曲:</strong> {chart.song?.arrangement?.name}</p>
          <p><strong>リリース時間:</strong> {chart.song?.releaseTime?.toDate().toLocaleString()}</p>
          <p><strong>原曲MV:</strong> <a href={chart.song?.originalVideo} target="_blank" rel="noopener noreferrer">{chart.song?.originalVideo}</a></p>
          <p><strong>ボーカルパターン:</strong></p>
          <ul>
            {chart.song?.vocalPatterns?.map((pattern, index) => (
              <li key={index}>
                {pattern.name} - {pattern.singers?.map(singer => singer.name).join(", ")}
              </li>
            ))}
          </ul>
          <p><strong>ユニット:</strong> {chart.song?.units?.map(unit => unit.name).join(", ")}</p>
        </div>
      </div>
    </div>
  );
};

export const ChartList = () => {
  const [charts, setCharts] = useState<Chart[]>([]);
  const [selectedChart, setSelectedChart] = useState<Chart | null>(null);

  useEffect(() => {
    const fetchCharts = async () => {
      const response = await masterClient.getCharts({});
      setCharts(response.charts);
    };

    fetchCharts().catch(console.error);
  }, []);

  return (
    <div className="chart-list-container">
      <h2>譜面一覧</h2>
      <div className="chart-list">
        {charts.map((chart) => (
          <div key={chart.id} className="chart-item">
            <div className="chart-thumbnail">
              {chart.song?.thumbnail ? (
                <img
                  src={
                    chart.song.thumbnail.startsWith("http")
                      ? chart.song.thumbnail
                      : `http://localhost:8888${chart.song.thumbnail}`
                  }
                  alt={chart.song.name}
                />
              ) : (
                <div className="no-thumbnail">No Image</div>
              )}
            </div>
            <div className="chart-info">
              <div className="chart-header">
                <h3 className="chart-name">{chart.song?.name}</h3>
                <div className="chart-difficulty">
                  {getDifficultyTypeDisplayName(chart.difficultyType)} {chart.level}
                </div>
              </div>
              <p className="chart-creators">
                {chart.song?.lyrics?.name} / {chart.song?.music?.name} / {chart.song?.arrangement?.name}
              </p>
            </div>
            <button
              className="detail-button"
              onClick={() => setSelectedChart(chart)}
            >
              詳細
            </button>
          </div>
        ))}
      </div>
      {selectedChart && (
        <ChartDetailModal
          chart={selectedChart}
          onClose={() => setSelectedChart(null)}
        />
      )}
    </div>
  );
};
