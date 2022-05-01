import React, { FC, useEffect, useState } from "react";
import { collection, getDocs, getFirestore, query, Timestamp, where } from "firebase/firestore";
import EventDropdown from "../components/eventdropdown";
import BookingItem from "../components/bookingitem";
import { Booking } from "../components/types";

const db = getFirestore();

const initialEvent = 20220506;
const initialGuests = [
  {
    BookingID: "#1651000440[20220506]Stefano",
    CustomerID: 0,
    CustomerName: "Stefano",
    EventID: 20220506,
    Timestamp: Timestamp.now(),
  }
];

const Bookings: FC = () => {
  const [event, setEvent] = useState(initialEvent);
  const [guests, setGuests] = useState(initialGuests);
  const bookedGuests = guests.map(
    b => <BookingItem
      key={b.BookingID}
      booking={b}
    />
  );

  useEffect(() => {
    const getGuests = async () => {
      const itemsColRef = query(collection(db, 'bookings'), where("EventID", "==", event));
      const listRef = await getDocs(itemsColRef);
      const list = listRef.docs;
      const g = list.map(bh => bh.data() as Booking);
      setGuests(g);
    }
    getGuests();
  }, [event]);

  const handleSelectEvent = (event: number) => {
    setEvent(event);
  }

  return (
    <div>
      <h3>{guests.length} guests for event
        <EventDropdown
          event={event}
          handleSelectEvent={handleSelectEvent}
        />
      </h3>
      <ol className="list-group list-group-numbered gap-1">{bookedGuests}</ol>
    </div>
  );
};

export default Bookings;