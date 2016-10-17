import React, { PropTypes } from 'react'

const Date = ({ onClick, date }) => (
    <li
        onClick={onClick}
        className="date"
    >
        { date }
    </li>
);

Date.propTypes = {
    onClick: PropTypes.func.isRequired,
    date: PropTypes.string.isRequired
}

export default Date