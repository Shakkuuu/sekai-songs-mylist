import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import "../styles/common.css";

export const HamburgerMenu = () => {
  const [isOpen, setIsOpen] = useState(false);

  const toggleMenu = () => {
    setIsOpen(!isOpen);
  };

  const handleClickOutside = (e: MouseEvent) => {
    const target = e.target as HTMLElement;
    if (isOpen && !target.closest('.nav-menu') && !target.closest('.hamburger-menu')) {
      setIsOpen(false);
    }
  };

  useEffect(() => {
    document.addEventListener('click', handleClickOutside);
    return () => {
      document.removeEventListener('click', handleClickOutside);
    };
  }, [isOpen]);

  return (
    <>
      <div className={`hamburger-menu ${isOpen ? "active" : ""}`} onClick={toggleMenu}>
        <div className="hamburger-icon"></div>
        <div className="hamburger-icon"></div>
        <div className="hamburger-icon"></div>
      </div>
      {isOpen && <div className="nav-overlay" />}
      <nav className={`nav-menu ${isOpen ? "active" : ""}`}>
        <ul>
          <li><Link to="/" onClick={toggleMenu}>トップ</Link></li>
          <li><Link to="/user" onClick={toggleMenu}>ユーザー</Link></li>
          <li><Link to="/signup" onClick={toggleMenu}>新規登録</Link></li>
          <li><Link to="/login" onClick={toggleMenu}>ログイン</Link></li>
          <li><Link to="/mylist" onClick={toggleMenu}>マイリスト</Link></li>
          <li><Link to="/master" onClick={toggleMenu}>Master</Link></li>
        </ul>
      </nav>
    </>
  );
};
