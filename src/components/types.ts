export interface Booking {
    BookingID: string;
    CustomerID: number;
    CustomerName: string;
    EventID: number;
    Timestamp: Date;
}

export interface Order {
    CategoryCode: string;
    CategoryName: string;
    CocktailCode: string;
    CocktailName: string;
    CustomerID: number;
    CustomerName: string;
    Done: boolean;
    OrderID: string;
    Quantity: number;
    Timestamp: Date;
}
