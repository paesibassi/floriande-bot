import React, { FC } from "react";

type Props = {
    cocktail: string,
    quantity: number,
}

const BatchItem: FC<Props> = (props: Props) => {
    const { cocktail, quantity } = props;
    return (
        <li className={`list-group-item border-2 rounded-3 p-1 fs-5`}>
            <span className="badge bg-primary">{quantity}</span> <strong>{cocktail}</strong>
        </li>
    );
}

export default BatchItem;