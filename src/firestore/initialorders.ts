import { Timestamp } from "firebase/firestore";

export const initialOrders = [
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