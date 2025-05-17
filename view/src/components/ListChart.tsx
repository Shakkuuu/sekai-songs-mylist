import { useEffect, useState } from "react";
import { masterClient } from "../lib/grpcClient";
import { Chart } from "../gen/master/chart_pb";
import { getDifficultyTypeDisplayName, getMusicVideoTypeDisplayName } from "../utils/enumDisplay";

export const ChartList = () => {
  const [charts, setSongs] = useState<Chart[]>([]);

  useEffect(() => {
    const fetchCharts = async () => {
      const response = await masterClient.getCharts({});
      setSongs(response.charts);
    };

    fetchCharts().catch(console.error);
  }, []);

  console.log(charts);

  return (
    <div>
      <h2>譜面一覧</h2>
      <ul>
        {charts.map((chart) => (
          <li key={chart.id}>
            <p><strong>ID:</strong> {chart.id}</p>
            <p><strong>Name:</strong> {chart.song?.name}</p>
            <p><strong>Lyrics:</strong> {chart.song?.lyrics?.name}</p>
            <p><strong>Music:</strong> {chart.song?.music?.name}</p>
            <p><strong>Arrangement:</strong> {chart.song?.arrangement?.name}</p>
            <p><strong>Thumbnail:</strong> {chart.song?.thumbnail}</p>
            <p><strong>Deleted:</strong> {chart.song?.deleted ? "Yes" : "No"}</p>
            <p><strong>Release Time:</strong> {chart.song?.releaseTime?.toDate().toLocaleString()}</p>
            <p><strong>Vocal Patterns:</strong></p>
            {chart.song?.vocalPatterns && chart.song?.vocalPatterns.length > 0 && (
              <ul>
                {chart.song?.vocalPatterns.map((pattern, index) => (
                  <li key={index}>
                    <p><strong>Name:</strong> {pattern.name}</p>
                    <p><strong>Singers:</strong> {pattern.singers?.map(singer => singer.name).join(", ")}</p>
                    <p><strong>Units:</strong> {pattern.units?.map(unit => unit.name).join(", ")}</p>
                  </li>
                ))}
              </ul>
            )}
            <p><strong>Music Video Types:</strong></p>
            <ul>
                {chart.song?.musicVideoTypes.map((type) =>
                    getMusicVideoTypeDisplayName(type)
                ).join(", ")}
            </ul>
            <p><strong>DifficultyType:</strong> {getDifficultyTypeDisplayName(chart.difficultyType)}</p>
            <p><strong>Level:</strong> {chart.level}</p>
            <p><strong>ChartViewLink:</strong> {chart.chartViewLink}</p>
          </li>
        ))}
      </ul>
    </div>
  );
};
