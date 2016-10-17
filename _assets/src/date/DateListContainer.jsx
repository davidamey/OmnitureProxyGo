import { connect } from 'react-redux'
import { selectDate } from './DateActions'
import DateList from './DateList'

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