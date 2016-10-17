// Component
import React, { PropTypes } from 'react'
import LogEntry from './LogEntry'

const Log = ({ logs }) => (
    <ul>
        { logs.map( (l, i) =>
            <LogEntry
                key={i}
                log={l}
            />
        )}
    </ul>
);

Log.propTypes = {
    logs : PropTypes.arrayOf(PropTypes.object.isRequired).isRequired
}

export { Log }

// Container
import { connect } from 'react-redux'

const mapStateToProps = (state) => {
    var visitors = state.logs[state.selectedDate] || null;

    return {
        logs : visitors ? visitors[state.selectedVisitor] || [] : []
    }
}

const LogContainer = connect(mapStateToProps)(Log)

export default LogContainer