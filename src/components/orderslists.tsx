import React, { FC } from "react";
import OrderItem from "./orderitem";
import { Order } from "./types";

type Props = {
  orders: Order[],
  handleServe: (OrderID: string) => void,
  handleDelete: (OrderID: string) => void,
}

const OrdersLists: FC<Props> = (props: Props) => {
  const { orders, handleServe, handleDelete } = props;
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
  return (
    <div>
      <h3>{outstandingOrders.length} outstanding orders</h3>
      <ol className="list-group list-group-numbered gap-1">{outstanding}</ol>
      <h3>{completedOrders.length} completed orders</h3>
      <ol className="list-group list-group-numbered gap-1">{completed}</ol>
    </div>
  );
};

export default OrdersLists;