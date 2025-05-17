import { useEffect, useState } from "react";
import { masterClient } from "../lib/grpcClient";
import { Singer } from "../gen/master/singer_pb";

export const SingerList = () => {
  const [singers, setSingers] = useState<Singer[]>([]);

  useEffect(() => {
    const fetchSingers = async () => {
      const response = await masterClient.getSingers({});
      setSingers(response.singers);
    };

    fetchSingers().catch(console.error);
  }, []);

  return (
    <div>
      <h2>シンガー一覧</h2>
      <ul>
        {singers.map((singer) => (
          <li key={singer.id}>{singer.id} : {singer.name}</li>
        ))}
      </ul>
    </div>
  );
};
