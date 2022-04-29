import React, { FC, useEffect, useState } from 'react';
import { initializeApp } from "firebase/app";
import { getFirestore, collection, onSnapshot, query, doc, updateDoc, getDocs, where, orderBy, deleteDoc, Timestamp } from "firebase/firestore";
import { getAuth, signInAnonymously } from "firebase/auth";
// TODO: Add SDKs for Firebase products that you want to use https://firebase.google.com/docs/web/setup#available-libraries
import OrdersLists from './components/orderslists';
import { Booking, Order } from "./components/types";
import EventDropdown from './components/eventdropdown';

const firebaseConfig = {
  apiKey: process.env.REACT_APP_APIKEY,
  authDomain: process.env.REACT_APP_AUTHDOMAIN,
  projectId: process.env.REACT_APP_PROJECTID,
};

const app = initializeApp(firebaseConfig);
const db = getFirestore();

export const authenticateAnonymously = () => {
  return signInAnonymously(getAuth(app));
};

async function handleServe(OrderID: string) {
  const orderRef = doc(db, "orders", OrderID);
  await updateDoc(orderRef, { Done: true });
}

async function handleDelete(OrderID: string) {
  const orderRef = doc(db, "orders", OrderID);
  await deleteDoc(orderRef);
}

const initialOrders = [
  {
    OrderID: "#1649442240Paolo",
    CustomerName: "Paolo", CustomerID: 0,
    Done: false,
    CategoryCode: "1Gin", CategoryName: "Gin", CocktailCode: "Dry Martini", CocktailName: "Dry Martini",
    Quantity: 1, Timestamp: new Date()
  },
  {
    OrderID: "#1649789999Silvano",
    CustomerName: "Silvano", CustomerID: 0,
    Done: false,
    CategoryCode: "1Rum", CategoryName: "Rum", CocktailCode: "2Daiquiri", CocktailName: "Daiquiri",
    Quantity: 1, Timestamp: new Date()
  },
  {
    OrderID: "#1650992794Federico",
    CustomerName: "Federico", CustomerID: 0,
    Done: true,
    CategoryCode: "1Gin", CategoryName: "Gin", CocktailCode: "2Negroni", CocktailName: "Negroni",
    Quantity: 1, Timestamp: new Date()
  },
];

const initialGuests = [
  {
    BookingID: "#1651000440[20220506]Stefano",
    CustomerID: 0,
    CustomerName: "Stefano",
    EventID: 20220506,
    Timestamp: new Date()
  }
];

const initialEvent = 20220506;

const App: FC = () => {
  const [orders, setOrders] = useState(initialOrders);
  const [event, setEvent] = useState(initialEvent);
  const [guests, setGuests] = useState(initialGuests);

  useEffect(() => {
    const getGuests = async () => {
      // const _ = await authenticateAnonymously(); 
      const itemsColRef = query(collection(db, 'bookings'), where("EventID", "==", event));
      const listRef = await getDocs(itemsColRef);
      const list = listRef.docs;
      const g = list.map(bh => bh.data() as Booking);
      setGuests(g);
    }
    getGuests();
  }, [event]);

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

  const handleSelectEvent = (event: number) => {
    setEvent(event);
  }

  return (
    <div className="container bg-light">
      <h1 className="display-1">Floriande Lounge</h1>
      <h2>
        {guests.length} guests for event
        <EventDropdown
          event={event}
          handleSelectEvent={handleSelectEvent}
        />
      </h2>
      <OrdersLists
        orders={orders}
        handleServe={handleServe}
        handleDelete={handleDelete}
      />
    </div>
  );
}

export default App;
