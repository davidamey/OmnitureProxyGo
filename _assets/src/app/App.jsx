import React from 'react'
import { connect } from 'react-redux'

import { fetchDates } from '../date/DateActions'
import { fetchVisitorsIfNeeded } from '../visitor/VisitorActions'
import { fetchLogIfNeeded } from '../log/LogActions'

import DateList from '../date/DateList'
import VisitorList from '../visitor/VisitorList'
import Log from '../log/Log'

class App extends React.Component {
    componentDidMount() {
        const { dispatch, selectedDate } = this.props
        dispatch(fetchDates())
    }

    componentWillReceiveProps(nextProps) {
        if (nextProps.selectedDate !== this.props.selectedDate) {
            const { dispatch, selectedDate } = nextProps
            dispatch(fetchVisitorsIfNeeded(selectedDate))
        }
        else if (nextProps.selectedVisitor !== this.props.selectedVisitor) {
            const { dispatch, selectedDate, selectedVisitor } = nextProps
            dispatch(fetchLogIfNeeded(selectedDate, selectedVisitor))
        }
    }

    render = () => (
        <div style={{ display: 'flex' }}>
            <DateList />
            <VisitorList />
            <Log />
        </div>
    )
};

function mapStateToProps(state) {
    return {
        selectedDate: state.selectedDate,
        selectedVisitor: state.selectedVisitor
    }
}

export default connect(mapStateToProps)(App)