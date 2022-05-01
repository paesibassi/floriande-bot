import React, { FC } from "react";
import { Booking } from "./types";

type Props = {
    booking: Booking,
}

const BookingItem: FC<Props> = (props: Props) => {
    const { booking } = props;
    const timeString = booking.Timestamp.toDate().toLocaleDateString('en-GB', {
        year: "numeric", month: "2-digit", day: "2-digit",
        hour: "2-digit", hour12: false, minute: "2-digit",
    });
    return (
        <li className={`list-group-item border-2 rounded-3 p-1`}>
            <strong>{booking.CustomerName}</strong><div className="d-none d-lg-inline"> booked on {timeString}</div><br/>
        </li>
    );
}

export default BookingItem;