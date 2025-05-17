import { useEffect, useState } from "react";
import { masterClient } from "../lib/grpcClient";
import { Unit } from "../gen/master/unit_pb";

export const UnitList = () => {
  const [units, setUnits] = useState<Unit[]>([]);

  useEffect(() => {
    const fetchUnits = async () => {
      const response = await masterClient.getUnits({});
      setUnits(response.units);
    };

    fetchUnits().catch(console.error);
  }, []);

  return (
    <div>
      <h2>ユニット一覧</h2>
      <ul>
        {units.map((unit) => (
          <li key={unit.id}>{unit.id} : {unit.name}</li>
        ))}
      </ul>
    </div>
  );
};
