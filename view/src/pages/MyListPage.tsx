import { useEffect, useState } from "react";
import { myListClient } from "../lib/grpcClient";
import { useNavigate } from "react-router-dom";
import {
  MyList,
  GetMyListsByUserIDRequest,
  ChangeMyListPositionRequest,
  DeleteMyListRequest,
  ChangeMyListNameRequest,
} from "../gen/mylist/v1/mylist_pb";

export const MyListPage = () => {
  const [myLists, setMyLists] = useState<MyList[]>([]);
  const [newName, setNewName] = useState("");
  const [dragIndex, setDragIndex] = useState<number | null>(null);
  const navigate = useNavigate();
  const [editId, setEditId] = useState<number | null>(null);
  const [editName, setEditName] = useState("");

  useEffect(() => {
    myListClient
      .getMyListsByUserID(new GetMyListsByUserIDRequest())
      .then((res) => setMyLists(res.myLists));
  }, []);

  const handleCreate = async () => {
    if (!newName) return;
    await myListClient.createMyList({
      name: newName,
      position: myLists.length + 1,
    });
    setNewName("");
    const res = await myListClient.getMyListsByUserID(
      new GetMyListsByUserIDRequest()
    );
    setMyLists(res.myLists);
  };

  // 並び替え
  const handleDragStart = (idx: number) => setDragIndex(idx);
  const handleDragOver = (e: React.DragEvent<HTMLLIElement>) =>
    e.preventDefault();
  const handleDrop = async (idx: number) => {
    if (dragIndex === null || dragIndex === idx) return;
    const newOrder = [...myLists];
    const [removed] = newOrder.splice(dragIndex, 1);
    newOrder.splice(idx, 0, removed);
    setMyLists(newOrder);
    setDragIndex(null);
    await myListClient.changeMyListPosition(
      new ChangeMyListPositionRequest({
        id: myLists.map((ml) => ml.id),
        position: myLists.map((_, idx) => idx + 1),
      })
    );
  };

  const handleEdit = (ml: MyList) => {
    setEditId(ml.id);
    setEditName(ml.name);
  };

  const handleEditSave = async () => {
    if (editId === null || !editName) return;
    await myListClient.changeMyListName(
      new ChangeMyListNameRequest({ id: editId, name: editName })
    );
    setEditId(null);
    setEditName("");
    const res = await myListClient.getMyListsByUserID(
      new GetMyListsByUserIDRequest()
    );
    setMyLists(res.myLists);
  };

  const handleEditCancel = () => {
    setEditId(null);
    setEditName("");
  };

  const handleDelete = async (id: number) => {
    await myListClient.deleteMyList(new DeleteMyListRequest({ id }));
    setMyLists(myLists.filter((ml) => ml.id !== id));
  };

  return (
    <div>
      <h2>マイリスト一覧</h2>
      <input
        value={newName}
        onChange={(e) => setNewName(e.target.value)}
        placeholder="新規マイリスト名"
      />
      <button onClick={handleCreate}>追加</button>
      <ul>
        {myLists
          .slice()
          .sort((a, b) => (a.position ?? 0) - (b.position ?? 0))
          .map((ml, idx) => (
            <li
              key={ml.id}
              draggable
              onDragStart={() => handleDragStart(idx)}
              onDragOver={handleDragOver}
              onDrop={() => handleDrop(idx)}
              style={{
                border: "1px solid #ccc",
                margin: "4px 0",
                padding: "4px",
                background: dragIndex === idx ? "#f0f0f0" : "white",
                cursor: "move",
              }}
            >
              {editId === ml.id ? (
                <>
                  <input
                    value={editName}
                    onChange={(e) => setEditName(e.target.value)}
                    style={{ marginRight: 8 }}
                  />
                  <button onClick={handleEditSave}>保存</button>
                  <button onClick={handleEditCancel}>キャンセル</button>
                </>
              ) : (
                <>
                  {ml.name}
                  <button onClick={() => navigate(`/mylist/${ml.id}`)}>
                    編集
                  </button>
                  <button onClick={() => handleEdit(ml)}>名前変更</button>
                  <button onClick={() => handleDelete(ml.id)}>削除</button>
                  <span style={{ marginLeft: 8, color: "#888" }}>⇅</span>
                </>
              )}
            </li>
          ))}
      </ul>
    </div>
  );
};
