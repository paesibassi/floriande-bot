import React, { FC, useEffect, useState } from "react";
import { collection, deleteDoc, doc, getFirestore, onSnapshot, orderBy, query, Timestamp, updateDoc, where } from "firebase/firestore";
import OrderItem from "../components/orderitem";
import { Order } from "../components/types";

const db = getFirestore();

const initialOrders = [
  {
    OrderID: "#1649442240Paolo",
    CustomerName: "Paolo", CustomerID: 0,
    Done: false,
    CategoryCode: "1Gin", CategoryName: "Gin", CocktailCode: "Dry Martini", CocktailName: "Dry Martini",
    Quantity: 1, Timestamp: Timestamp.now(),
  },
  {
    OrderID: "#1650992794Federico",
    CustomerName: "Federico", CustomerID: 0,
    Done: true,
    CategoryCode: "1Gin", CategoryName: "Gin", CocktailCode: "2Negroni", CocktailName: "Negroni",
    Quantity: 1, Timestamp: Timestamp.now(),
  },
];

async function handleServe(OrderID: string) {
  const orderRef = doc(db, "orders", OrderID);
  await updateDoc(orderRef, { Done: true });
}

async function handleDelete(OrderID: string) {
  const orderRef = doc(db, "orders", OrderID);
  await deleteDoc(orderRef);
}

const Orders: FC = () => {
  const [orders, setOrders] = useState(initialOrders);
  const outstandingOrders = orders.filter(o => o.Done === false)
  const completedOrders = orders.filter(o => o.Done === true)
  const outstanding = outstandingOrders.map(
    o => <OrderItem
      key={o.OrderID}
      order={o}
      handleServe={handleServe}
      handleDelete={handleDelete}
    />
  );
  const completed = completedOrders.map(
    o => <OrderItem
      key={o.OrderID}
      order={o}
      handleDelete={handleDelete}
    />
  );

  useEffect(() => {
    const listenToOrders = () => {
      const last24Hours = 7 * 24 * 60 * 60 * 1000
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
    <div>
      <h3>{outstandingOrders.length} outstanding orders</h3>
      <ol className="list-group list-group-numbered gap-1">{outstanding}</ol>
      <h3>{completedOrders.length} completed orders</h3>
      <ol className="list-group list-group-numbered gap-1">{completed}</ol>
    </div>
  );
};

export default Orders;