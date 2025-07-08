import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { myListClient } from "../lib/grpcClient";
import {
  MyList,
  GetMyListsByUserIDRequest,
  ChangeMyListPositionRequest,
  ChangeMyListNameRequest,
  DeleteMyListRequest,
} from "../gen/mylist/v1/mylist_pb";
import { HamburgerMenu } from "../components/HamburgerMenu";
import "../styles/common.css";
import "./MyListPage.css";

export const MyListPage = () => {
  const [myLists, setMyLists] = useState<MyList[]>([]);
  const [newName, setNewName] = useState("");
  const [dragIndex, setDragIndex] = useState<number | null>(null);
  const navigate = useNavigate();
  const [editId, setEditId] = useState<number | null>(null);
  const [editName, setEditName] = useState("");

  useEffect(() => {
    const fetchMyLists = async () => {
      const response = await myListClient.getMyListsByUserID(
        new GetMyListsByUserIDRequest()
      );
      setMyLists(response.myLists);
    };
    fetchMyLists();
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
  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) =>
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
    <div className="container">
      <HamburgerMenu />
      <div className="page-header">
        <h1>マイリスト</h1>
      </div>
      <div className="card">
        <div className="form-group">
          <input
            type="text"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            placeholder="新しいマイリスト名"
          />
          <button className="button" onClick={handleCreate}>
            新規マイリスト作成
          </button>
        </div>
      </div>
      <div className="mylist-grid">
        {myLists.map((myList) => (
          <div
            key={myList.id}
            className="card mylist-card"
            draggable
            onDragStart={() => handleDragStart(myLists.indexOf(myList))}
            onDragOver={handleDragOver}
            onDrop={() => handleDrop(myLists.indexOf(myList))}
          >
            {editId === myList.id ? (
              <div className="edit-form">
                <input
                  type="text"
                  value={editName}
                  onChange={(e) => setEditName(e.target.value)}
                />
                <button onClick={handleEditSave}>保存</button>
                <button onClick={handleEditCancel}>キャンセル</button>
              </div>
            ) : (
              <>
                <h3 onClick={() => navigate(`/mylist/${myList.id}`)}>
                  {myList.name}
                </h3>
                <p>作成日: {myList.createdAt?.toDate().toLocaleDateString()}</p>
                <div className="mylist-actions">
                  <button onClick={() => handleEdit(myList)}>
                    マイリスト名編集
                  </button>
                  <button onClick={() => handleDelete(myList.id)}>削除</button>
                </div>
              </>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};
