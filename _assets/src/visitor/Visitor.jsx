import React, { PropTypes } from 'react'

const Visitor = ({ onClick, vid }) => (
    <li
        onClick={onClick}
        className="visitor"
    >
        { vid }
    </li>
);

Visitor.propTypes = {
    onClick: PropTypes.func.isRequired,
    vid: PropTypes.string.isRequired
}

export default Visitor