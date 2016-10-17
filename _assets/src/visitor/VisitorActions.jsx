import fetch from 'isomorphic-fetch'

export const SELECT_VISITOR = 'SELECT_VISITOR'
export const REQUEST_VISITORS = 'REQUEST_VISITORS'
export const RECEIVE_VISITORS = 'RECEIVE_VISITORS'

function requestVisitors(date) {
    return {
       type: REQUEST_VISITORS,
       date
    }
}

function receiveVisitors(date, json) {
    return {
        type: RECEIVE_VISITORS,
        date,
        visitors: json,
        receivedAt: Date.now()
    }
}

function fetchVisitors(date) {
    return dispatch => {
        dispatch(requestVisitors(date))
        return fetch(`/api/logs/${date}/`)
            .then(response => response.json())
            .then(json => dispatch(receiveVisitors(date, json)))
    }
}

export function fetchVisitorsIfNeeded(date) {
    return dispatch => {
        if (date != null) {
            return dispatch(fetchVisitors(date))
        }
    }
}

export function selectVisitor(vid) {
    return {
        type: SELECT_VISITOR,
        vid
    }
}