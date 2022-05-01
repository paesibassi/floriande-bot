import React, { FC, useEffect } from 'react';
import { Link, Outlet, useLocation } from 'react-router-dom';
import { authenticate } from './firestore/db';

const App: FC = () => {
  const location = useLocation();
  const pathname = location.pathname;

  useEffect(() => {
    authenticate();
  }, []);

  return (
    <div className="container bg-light">
      <h1 className="display-1">Floriande Lounge</h1>
      <ul className="nav nav-tabs">
        <li className="nav-item">
          <Link className={pathname === "/" ? "nav-link active" : "nav-link"} to="/">Orders</Link>
        </li>
        <li className="nav-item">
          <Link className={pathname === "/bookings" ? "nav-link active" : "nav-link"} to="/bookings">Bookings</Link>
        </li>
      </ul>
      <Outlet />
    </div>
  );
}

export default App;
