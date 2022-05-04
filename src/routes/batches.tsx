import React, { FC } from "react";
import { useOrders } from "../App";
import BatchItem from "../components/batchitem";
import { Order } from "../components/types";

type ordersGroups = {
  [key: string]: number;
};

const Orders: FC = () => {
  const { orders } = useOrders();
  const ordersGrouped: ordersGroups = orders
    .filter(o => o.Done === false)
    .reduce(
      (counts: ordersGroups, order: Order) => {
        const cocktail = order.CocktailName;
        const count = counts[order.CocktailName] || 0;
        return { ...counts, [cocktail]: count + 1 }
      },
      {},
    );
  const batches = Object.keys(ordersGrouped).map(
    c => <BatchItem
      key={c}
      cocktail = {c}
      quantity={ordersGrouped[c]}
    />
  );

  return (
    <div>
      <ul className="list-group gap-1">{batches}</ul>
    </div>
  );
};

export default Orders;