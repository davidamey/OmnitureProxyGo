// Component
import React, { PropTypes } from 'react'
import Date from './Date'

const DateList = ({ dates, onDateClick }) => (
    <ul>
        { dates.map( d =>
            <Date
                key={d}
                date={d}
                onClick={ () => onDateClick(d) }
            />
        )}
    </ul>
);

DateList.propTypes = {
    dates : PropTypes.arrayOf(PropTypes.string.isRequired).isRequired,
    onDateClick: PropTypes.func.isRequired
}

export { DateList }

// Container
import { connect } from 'react-redux'
import { selectDate } from './DateActions'

const mapStateToProps = (state) => {
    return {
        dates: Object.keys(state.logs)
    }
}

const mapDispatchToProps = (dispatch) => {
    return {
        onDateClick: (date) => {
            // alert(`date ${id} clicked`);
            dispatch(selectDate(date))
        }
    }
}

const DateListContainer = connect(
    mapStateToProps,
    mapDispatchToProps
)(DateList)

export default DateListContainer