export const SELECT_DATE = 'SELECT_DATE'
export const REQUEST_DATES = 'REQUEST_DATES'
export const RECEIVE_DATES = 'RECEIVE_DATES'

function requestDates(date) {
    return {
       type: REQUEST_DATES
    }
}

function receiveDates(json) {
    return {
        type: RECEIVE_DATES,
        dates: json,
        receivedAt: Date.now()
    }
}

export function fetchDates() {
    return dispatch => {
        dispatch(requestDates())
        return fetch(`/api/logs/`)
            .then(response => response.json())
            .then(json => dispatch(receiveDates(json)))
    }
}

export function selectDate(date) {
    return {
        type: SELECT_DATE,
        date
    }
}