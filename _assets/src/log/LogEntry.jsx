import React, { PropTypes } from 'react'

const LogEntry = ({ log }) => (
    <li className="LogEntry">
        { log.pageName }
    </li>
);

LogEntry.propTypes = {
    log: PropTypes.object.isRequired
}

export default LogEntry