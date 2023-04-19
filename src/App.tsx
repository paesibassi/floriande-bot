import React, { FC, useEffect, useState } from 'react';
import { Link, Outlet, useLocation, useOutletContext } from 'react-router-dom';
import {
  collection, deleteDoc, doc, getFirestore,
  onSnapshot, orderBy, query, Timestamp, where
} from 'firebase/firestore';
import { authenticate } from './firestore/db';
import { Order } from './components/types';

const db = getFirestore();
const floriandeAPI = process.env.REACT_APP_APIENDPOINT;

type ContextType = {
  orders: Order[],
  handleServe: (OrderID: string) => void,
  handleDelete: (OrderID: string) => void,
}

async function handleServe(OrderID: string) {
  const command = { "command": "CloseOrder", "orderID": OrderID };
  if (floriandeAPI === undefined) { 
    console.error("Could not get the API endpoint address from env")
    return
  }
  const response = await fetch(floriandeAPI, {
    method: 'POST',
    mode: 'no-cors',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(command),
  });
  if (!response.ok && response.status !== 0) {
    console.error(`An error has occured: ${response.status}`)
  }
}

async function handleDelete(OrderID: string) {
  const orderRef = doc(db, "orders", OrderID);
  await deleteDoc(orderRef);
}

const App: FC = () => {
  const [orders, setOrders] = useState([] as Order[]);
  const location = useLocation();
  const pathname = location.pathname;

  useEffect(() => {
    authenticate();
  }, []);

  useEffect(() => {
    const listenToOrders = () => {
      const last24Hours = 24 * 60 * 60 * 1000
      const timestamp = Timestamp.fromDate(new Date(Date.now() - last24Hours))
      const q = query(collection(db, "orders"), where("Timestamp", ">=", timestamp), orderBy("Timestamp"));
      const unsubscribe = onSnapshot(q, (querySnapshot) => {
        const o = querySnapshot.docs.map(bh => bh.data() as Order);
        setOrders(o);
      });
      return unsubscribe;
    };
    listenToOrders();
  }, []);

  return (
    <div className="container bg-light">
      <h1 className="display-1">Floriande Lounge</h1>
      <ul className="nav nav-tabs">
        <li className="nav-item">
          <Link className={pathname === "/" ? "nav-link active" : "nav-link"} to="/">Orders</Link>
        </li>
        <li className="nav-item">
          <Link className={pathname === "/batches" ? "nav-link active" : "nav-link"} to="/batches">Batches</Link>
        </li>
        <li className="nav-item">
          <Link className={pathname === "/bookings" ? "nav-link active" : "nav-link"} to="/bookings">Bookings</Link>
        </li>
      </ul>
      <Outlet context={{ orders, handleServe, handleDelete }}/>
    </div>
  );
}

export function useOrders() {
  return useOutletContext<ContextType>();
}

export default App;
