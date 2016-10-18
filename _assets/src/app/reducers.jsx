import { SELECT_DATE, REQUEST_DATES, RECEIVE_DATES } from '../date/DateActions'
import { SELECT_VISITOR, REQUEST_VISITORS, RECEIVE_VISITORS } from '../visitor/VisitorActions'
import { REQUEST_LOG, RECEIVE_LOG } from '../log/LogActions'

const initialState = {
    logs: {
        /*
        '2016-10-14': {
            'vis1': ['v1-log1', 'v1-log2'],
            'vis2': ['v2-log1', 'v2-log2']
        }
        */
    },

    selectedDate : null,
    selectedVisitor : null,

    isFetchingDates : false,
    isFetchingVisitors : false,
    isFetchingLog : false
};

const reducers = (state = initialState, action) => {
    switch (action.type) {
        case SELECT_DATE:
            return { ...state, selectedDate: action.date }

        case REQUEST_DATES:
            return { ...state, isFetchingDates: true }

        case RECEIVE_DATES:
            var dates = {};
            action.dates.forEach(d => {
                dates[d] = {}
            });

            return {
                ...state,
                isFetchingDates: false,
                logs: dates
            }

        case SELECT_VISITOR:
            return { ...state, selectedVisitor: action.vid }

        case REQUEST_VISITORS:
            return { ...state, isFetchingVisitors: true }

        case RECEIVE_VISITORS:
            var visitors = {};
            action.visitors.forEach(v => {
                visitors[v] = []
            });

            return {
                ...state,
                isFetchingVisitors: false,
                logs: {
                    ...state.logs,
                    [action.date]: visitors
                }
            }

        case REQUEST_LOG:
            return { ...state, isFetchingLog: true }

        case RECEIVE_LOG:
            return {
                ...state,
                isFetchingLog: false,
                logs: {
                    ...state.logs,
                    [action.date]: {
                        ...state.logs[action.date],
                        [action.visitor]: action.log
                    }
                }
            }

        default:
            return state   
    }
}

export default reducers