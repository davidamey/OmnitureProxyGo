import React, { PropTypes } from 'react'

const Visitor = ({ onClick, vid, selected }) => (
    <li
        onClick={onClick}
        className="visitor"
        style={{
            background: selected ? 'lightgreen' : 'auto'
        }}
    >
        { vid }
    </li>
);

Visitor.propTypes = {
    onClick: PropTypes.func.isRequired,
    vid: PropTypes.string.isRequired,
    selected: PropTypes.bool
}

export default Visitor