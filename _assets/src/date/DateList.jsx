// Component
import React, { PropTypes } from 'react'
import DateItem from './DateItem'

const DateList = ({ dates, selectedDate, onDateClick }) => (
    <ul>
        { dates.map( d =>
            <DateItem
                key={d}
                date={d}
                selected={ d === selectedDate }
                onClick={ () => onDateClick(d) }
            />
        )}
    </ul>
);

DateList.propTypes = {
    dates: PropTypes.arrayOf(PropTypes.string.isRequired).isRequired,
    selectedDate: PropTypes.string,
    onDateClick: PropTypes.func.isRequired
}

export { DateList }

// Container
import { connect } from 'react-redux'
import { selectDate } from './DateActions'

const mapStateToProps = (state) => {
    return {
        dates: Object.keys(state.logs),
        selectedDate: state.selectedDate
    }
}

const mapDispatchToProps = (dispatch) => {
    return {
        onDateClick: (date) => {
            dispatch(selectDate(date))
        }
    }
}

const DateListContainer = connect(
    mapStateToProps,
    mapDispatchToProps
)(DateList)

export default DateListContainer