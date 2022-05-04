import { Timestamp } from "firebase/firestore";

export const initialOrders = [
  {
    OrderID: "#1649442240Paolo",
    CustomerName: "Paolo", CustomerID: 0, CustomerLang: "it",
    Done: false,
    CategoryCode: "1Gin", CategoryName: "Gin", CocktailCode: "Dry Martini", CocktailName: "Dry Martini",
    Quantity: 1, Timestamp: Timestamp.now(),
  },
  {
    OrderID: "#1649789999Silvano",
    CustomerName: "Silvano", CustomerID: 1, CustomerLang: "it",
    Done: false,
    CategoryCode: "1Rum", CategoryName: "Rum", CocktailCode: "2Cuba Libre", CocktailName: "Cuba Libre",
    Quantity: 1, Timestamp: Timestamp.now(),
  },
  {
    OrderID: "#1651093312Alessio",
    CustomerName: "Alessio", CustomerID: 2, CustomerLang: "it",
    Done: false,
    CategoryCode: "1Rum", CategoryName: "Rum", CocktailCode: "2Cuba Libre", CocktailName: "Cuba Libre",
    Quantity: 1, Timestamp: Timestamp.now(),
  },
  {
    OrderID: "#1650992794Federico",
    CustomerName: "Federico", CustomerID: 3, CustomerLang: "en",
    Done: true,
    CategoryCode: "1Gin", CategoryName: "Gin", CocktailCode: "2Negroni", CocktailName: "Negroni",
    Quantity: 1, Timestamp: Timestamp.now(),
  },
];