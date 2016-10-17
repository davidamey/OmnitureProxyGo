import fetch from 'isomorphic-fetch'

export const REQUEST_LOG = 'REQUEST_LOG'
export const RECEIVE_LOG = 'RECEIVE_LOG'

function requestLog(date, visitor) {
    return {
       type: REQUEST_LOG,
       date,
       visitor
    }
}

function receiveLog(date, visitor, json) {
    return {
        type: RECEIVE_LOG,
        date,
        visitor,
        log: json,
        receivedAt: Date.now()
    }
}

function fetchLog(date, visitor) {
    return dispatch => {
        dispatch(requestLog(date, visitor))
        return fetch(`/api/logs/${date}/${visitor}`)
            .then(response => response.json())
            .then(json => dispatch(receiveLog(date, visitor, json)))
    }
}

export function fetchLogIfNeeded(date, visitor) {
    return dispatch => {
        if (date != null && visitor != null) {
            return dispatch(fetchLog(date, visitor))
        }
    }
}