import React, { FC } from "react";
import { Order } from "./types";

type Props = {
    order: Order,
    handleServe?: (OrderID: string) => void,
    handleDelete: (OrderID: string) => void,
}

const OrderItem: FC<Props> = (props: Props) => {
    const { order, handleServe, handleDelete } = props;
    const date = order.Timestamp.toDate();
    const timeString = `${date.getHours().toString().padStart(2, "0")}:${date.getMinutes().toString().padStart(2, "0")}`;
    return (
        <li className={`list-group-item border-2 rounded-3 p-1 ${order.Done ? "bg-secondary bg-gradient" : ""}`}>
            <strong>{order.CocktailName}</strong> &gt; {order.CustomerName} <div className="d-none d-md-inline">@ {timeString}</div><br/>
            <small>{order.OrderID}</small>
            <div className="position-absolute top-50 end-0 translate-middle-y">
                {handleServe &&
                    <button type="button" className="btn btn-primary mx-1" onClick={() => handleServe(order.OrderID)}>
                        <i className="bi bi-cup-straw"></i>
                    </button>
                }
                <button type="button" className="btn btn-danger mx-1" onClick={() => handleDelete(order.OrderID)}>
                    <i className="bi bi-x-circle"></i>
                </button>
            </div>
        </li>
    );
}

export default OrderItem;