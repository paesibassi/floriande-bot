import React, { FC, useEffect, useState } from "react";
import { collection, getDocs, getFirestore, query, where } from "firebase/firestore";
import EventDropdown from "../components/eventdropdown";
import BookingItem from "../components/bookingitem";
import { Booking } from "../components/types";

const db = getFirestore();

const defaultEvent = 20220506;

const Bookings: FC = () => {
  const [event, setEvent] = useState(defaultEvent);
  const [bookings, setBookings] = useState([] as Booking[]);
  const bookingItems = bookings.map(
    b => <BookingItem
      key={b.BookingID}
      booking={b}
    />
  );

  useEffect(() => {
    const getBookings = async () => {
      const itemsColRef = query(collection(db, 'bookings'), where("EventID", "==", event));
      const listRef = await getDocs(itemsColRef);
      const list = listRef.docs;
      const g = list.map(bh => bh.data() as Booking);
      setBookings(g);
    }
    getBookings();
  }, [event]);

  const handleSelectEvent = (event: number) => {
    setEvent(event);
  }

  return (
    <div>
      <h3>{bookings.length} guests for event
        <EventDropdown
          event={event}
          handleSelectEvent={handleSelectEvent}
        />
      </h3>
      <ol className="list-group list-group-numbered gap-1">{bookingItems}</ol>
    </div>
  );
};

export default Bookings;