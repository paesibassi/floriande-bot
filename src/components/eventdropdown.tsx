import React, { FC } from "react";

type Props = {
  event: number;
  handleSelectEvent: (event: number) => void,
}

const EventDropdown: FC<Props> = (props: Props) => {
  const { event, handleSelectEvent } = props;
  const events = [20220506, 20220815];
  return (
    <div className="dropdown float-end">
      <button
        className="btn btn-primary btn-sm dropdown-toggle"
        type="button"
        id="dropdownEvent"
        data-bs-toggle="dropdown"
        aria-expanded="false"
      >
        { event }
      </button>
      <ul className="dropdown-menu" aria-labelledby="dropdownEvent">
        {events.map(e => 
          <li key={e}><button className="dropdown-item" type="button" onClick={() => handleSelectEvent(e)}>{e}</button></li>
        )}
      </ul>
    </div>
  );
}

export default EventDropdown;
